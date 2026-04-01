package integration_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"apiRest/internal/entity"
	"apiRest/internal/infra/auth"
	deliveryhttp "apiRest/internal/infra/delivery/http"
	"apiRest/internal/infra/messaging"
	"apiRest/internal/infra/observability"
	"apiRest/internal/infra/worker"
	"apiRest/internal/usecase"

	"github.com/gin-gonic/gin"
)

type integrationRepo struct {
	items map[string]*entity.Notification
}

func (r *integrationRepo) Save(ctx context.Context, n *entity.Notification) error {
	if r.items == nil {
		r.items = map[string]*entity.Notification{}
	}
	clone := *n
	r.items[n.ID] = &clone
	return nil
}

func (r *integrationRepo) UpdateStatus(ctx context.Context, id, status string) error {
	if r.items == nil {
		return nil
	}
	if n, ok := r.items[id]; ok {
		n.Status = status
	}
	return nil
}

type okSender struct{}

func (s *okSender) Send(ctx context.Context, n entity.Notification) error {
	return nil
}

func TestNotificationFlowProtectedRouteAndAsyncProcessing(t *testing.T) {
	gin.SetMode(gin.TestMode)

	metrics := observability.NewMetrics()
	repo := &integrationRepo{}
	queue := messaging.NewLocalQueue(10)
	uc := &usecase.SendNotificationUseCase{Repo: repo, Publisher: queue}
	handler := &deliveryhttp.NotificationHandler{UseCase: uc, Metrics: metrics}

	jwtManager := auth.NewJWTManager("test-secret", "test", time.Hour)
	token, err := jwtManager.Generate("user-123")
	if err != nil {
		t.Fatalf("token generation failed: %v", err)
	}

	w := &worker.NotificationWorker{
		Queue:       queue,
		Repo:        repo,
		Sender:      &okSender{},
		Metrics:     metrics,
		MaxRetries:  3,
		BaseBackoff: time.Millisecond,
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	if err := w.Start(ctx); err != nil {
		t.Fatalf("worker start failed: %v", err)
	}

	router := gin.New()
	protected := router.Group("/")
	protected.Use(deliveryhttp.JWTMiddleware(jwtManager))
	protected.POST("/notifications", handler.Create)

	payload, _ := json.Marshal(map[string]string{
		"content": "pagamento aprovado",
		"type":    "payment",
	})

	req := httptest.NewRequest(http.MethodPost, "/notifications", bytes.NewReader(payload))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	res := httptest.NewRecorder()

	router.ServeHTTP(res, req)

	if res.Code != http.StatusCreated {
		t.Fatalf("expected status 201, got %d", res.Code)
	}

	deadline := time.Now().Add(500 * time.Millisecond)
	for time.Now().Before(deadline) {
		for _, n := range repo.items {
			if n.Status == entity.StatusSent {
				return
			}
		}
		time.Sleep(10 * time.Millisecond)
	}

	t.Fatal("expected one notification processed with sent status")
}
