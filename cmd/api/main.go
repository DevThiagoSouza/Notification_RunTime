package main

import (
	"context"
	"database/sql"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"apiRest/internal/infra/auth"
	httpdelivery "apiRest/internal/infra/delivery/http"
	"apiRest/internal/infra/messaging"
	"apiRest/internal/infra/observability"
	"apiRest/internal/infra/realtime"
	"apiRest/internal/infra/repository"
	"apiRest/internal/infra/worker"
	"apiRest/internal/usecase"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	metrics := observability.NewMetrics()

	dbURL := getEnv("DATABASE_URL", "postgres://postgres:postgres@localhost:5432/notifications?sslmode=disable")
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal("Erro ao conectar ao banco de dados: ", err)
	}
	if err := db.PingContext(ctx); err != nil {
		log.Fatal("Erro ao validar conexão com banco de dados: ", err)
	}
	defer db.Close()

	hub := realtime.NewHub()
	notifier := &realtime.Notifier{Hub: hub, Logger: logger}

	queue, err := buildQueue(logger)
	if err != nil {
		log.Fatal("Erro ao configurar fila: ", err)
	}
	defer queue.Close()

	repo := repository.NewPostgresNotificationRepository(db)
	useCase := &usecase.SendNotificationUseCase{Repo: repo, Publisher: queue}
	handler := &httpdelivery.NotificationHandler{UseCase: useCase, Metrics: metrics}

	workerService := &worker.NotificationWorker{
		Queue:       queue,
		Repo:        repo,
		Sender:      notifier,
		Logger:      logger,
		Metrics:     metrics,
		MaxRetries:  3,
		BaseBackoff: time.Second,
	}
	if err := workerService.Start(ctx); err != nil {
		log.Fatal("Erro ao iniciar worker: ", err)
	}

	jwtSecret := getEnv("JWT_SECRET", "dev-secret-change-me")
	if jwtSecret == "" {
		jwtSecret = "dev-secret-change-me"
	}
	jwtManager := auth.NewJWTManager(jwtSecret, "apiRest", 24*time.Hour)
	authHandler := &httpdelivery.AuthHandler{JWTManager: jwtManager}
	wsHandler := &httpdelivery.WebSocketHandler{Hub: hub, Logger: logger}
	metricsHandler := &httpdelivery.MetricsHandler{Metrics: metrics}

	router := gin.Default()

	router.POST("/auth/token", authHandler.Token)
	router.GET("/metrics", metricsHandler.Get)

	protected := router.Group("/")
	protected.Use(httpdelivery.JWTMiddleware(jwtManager))
	protected.POST("/notifications", handler.Create)
	protected.GET("/ws", wsHandler.Handle)

	port := getEnv("PORT", "8080")
	log.Printf("Servidor rodando na porta %s", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatal("Erro ao iniciar o servidor: ", err)
	}

}

func buildQueue(logger *slog.Logger) (messaging.Queue, error) {
	rabbitURL := os.Getenv("RABBITMQ_URL")
	if rabbitURL == "" {
		if logger != nil {
			logger.Info("using local in-memory queue")
		}
		return messaging.NewLocalQueue(256), nil
	}

	queueName := os.Getenv("RABBITMQ_QUEUE")
	queue, err := messaging.NewRabbitMQQueue(rabbitURL, queueName)
	if err != nil {
		return nil, err
	}
	if logger != nil {
		logger.Info("using rabbitmq queue", "queue", queueName)
	}
	return queue, nil
}

func getEnv(key, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	return value
}
