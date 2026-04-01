package worker

import (
	"apiRest/internal/entity"
	"apiRest/internal/infra/messaging"
	"apiRest/internal/infra/observability"
	"context"
	"log/slog"
	"time"
)

type NotificationStatusRepository interface {
	UpdateStatus(ctx context.Context, notificationID, status string) error
}

type NotificationSender interface {
	Send(ctx context.Context, notification entity.Notification) error
}

type NotificationWorker struct {
	Queue       messaging.Queue
	Repo        NotificationStatusRepository
	Sender      NotificationSender
	Logger      *slog.Logger
	Metrics     *observability.Metrics
	MaxRetries  int
	BaseBackoff time.Duration
}

func (w *NotificationWorker) Start(ctx context.Context) error {
	if w.MaxRetries <= 0 {
		w.MaxRetries = 3
	}
	if w.BaseBackoff <= 0 {
		w.BaseBackoff = time.Second
	}

	return w.Queue.Consume(ctx, w.handle)
}

func (w *NotificationWorker) handle(ctx context.Context, notification entity.Notification) error {
	if w.Logger != nil {
		w.Logger.Info("processing notification", "notification_id", notification.ID, "user_id", notification.UserID)
	}
	if w.Metrics != nil {
		w.Metrics.IncProcessed()
	}

	if err := w.Repo.UpdateStatus(ctx, notification.ID, entity.StatusProcessing); err != nil && w.Logger != nil {
		w.Logger.Error("failed to set status processing", "notification_id", notification.ID, "error", err)
	}

	var lastErr error
	for attempt := 1; attempt <= w.MaxRetries; attempt++ {
		lastErr = w.Sender.Send(ctx, notification)
		if lastErr == nil {
			if err := w.Repo.UpdateStatus(ctx, notification.ID, entity.StatusSent); err != nil && w.Logger != nil {
				w.Logger.Error("failed to set status sent", "notification_id", notification.ID, "error", err)
			}
			if w.Metrics != nil {
				w.Metrics.IncSent()
			}
			if w.Logger != nil {
				w.Logger.Info("notification sent", "notification_id", notification.ID, "attempt", attempt)
			}
			return nil
		}

		if w.Metrics != nil {
			w.Metrics.IncRetried()
		}
		if w.Logger != nil {
			w.Logger.Warn("notification send failed, retrying", "notification_id", notification.ID, "attempt", attempt, "error", lastErr)
		}

		backoff := w.BaseBackoff * time.Duration(1<<(attempt-1))
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(backoff):
		}
	}

	if err := w.Repo.UpdateStatus(ctx, notification.ID, entity.StatusFailed); err != nil && w.Logger != nil {
		w.Logger.Error("failed to set status failed", "notification_id", notification.ID, "error", err)
	}
	if w.Metrics != nil {
		w.Metrics.IncFailed()
	}
	if w.Logger != nil {
		w.Logger.Error("notification failed after max retries", "notification_id", notification.ID, "error", lastErr)
	}

	return lastErr
}
