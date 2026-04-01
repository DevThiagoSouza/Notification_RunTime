package worker_test

import (
	"context"
	"errors"
	"sync"
	"testing"
	"time"

	"apiRest/internal/entity"
	"apiRest/internal/infra/messaging"
	"apiRest/internal/infra/observability"
	"apiRest/internal/infra/worker"
)

type inMemoryRepo struct {
	mu       sync.Mutex
	statuses map[string]string
}

func (r *inMemoryRepo) UpdateStatus(ctx context.Context, id, status string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.statuses == nil {
		r.statuses = map[string]string{}
	}
	r.statuses[id] = status
	return nil
}

func (r *inMemoryRepo) status(id string) string {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.statuses[id]
}

type flakySender struct {
	mu       sync.Mutex
	fails    int
	attempts int
}

func (s *flakySender) Send(ctx context.Context, n entity.Notification) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.attempts++
	if s.attempts <= s.fails {
		return errors.New("temporary send error")
	}
	return nil
}

func TestNotificationWorkerRetryAndSend(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	queue := messaging.NewLocalQueue(10)
	repo := &inMemoryRepo{}
	sender := &flakySender{fails: 2}
	metrics := observability.NewMetrics()

	w := &worker.NotificationWorker{
		Queue:       queue,
		Repo:        repo,
		Sender:      sender,
		Metrics:     metrics,
		MaxRetries:  3,
		BaseBackoff: 10 * time.Millisecond,
	}

	if err := w.Start(ctx); err != nil {
		t.Fatalf("start failed: %v", err)
	}

	n := &entity.Notification{ID: "n1", UserID: "u1", Content: "hello", Type: "comment", Status: entity.StatusPending}
	if err := queue.Publish(ctx, n); err != nil {
		t.Fatalf("publish failed: %v", err)
	}

	deadline := time.Now().Add(500 * time.Millisecond)
	for time.Now().Before(deadline) {
		if got := repo.status("n1"); got == entity.StatusSent {
			break
		}
		time.Sleep(10 * time.Millisecond)
	}

	if got := repo.status("n1"); got != entity.StatusSent {
		t.Fatalf("expected status sent, got %s", got)
	}

	if sender.attempts != 3 {
		t.Fatalf("expected 3 attempts, got %d", sender.attempts)
	}

	snapshot := metrics.Snapshot()
	if snapshot["retried"] < 2 {
		t.Fatalf("expected at least 2 retries, got %d", snapshot["retried"])
	}
}
