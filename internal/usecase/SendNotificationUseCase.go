package usecase

import (
	"context"

	"apiRest/internal/entity"
)

type NotificationRepository interface {
	Save(ctx context.Context, notification *entity.Notification) error
}

type NotificationPublisher interface {
	Publish(ctx context.Context, notification *entity.Notification) error
}

type SendNotificationUseCase struct {
	Repo      NotificationRepository
	Publisher NotificationPublisher
}

func (uc *SendNotificationUseCase) Execute(ctx context.Context, notification *entity.Notification) error {
	if err := notification.Validate(); err != nil {
		return err
	}

	notification.Status = entity.StatusPending
	if err := uc.Repo.Save(ctx, notification); err != nil {
		return err
	}

	if uc.Publisher == nil {
		return nil
	}

	return uc.Publisher.Publish(ctx, notification)
}
