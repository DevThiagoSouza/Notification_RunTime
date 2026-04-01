package messaging

import (
	"context"

	"apiRest/internal/entity"
)

type Queue interface {
	Publish(ctx context.Context, notification *entity.Notification) error
	Consume(ctx context.Context, handler func(context.Context, entity.Notification) error) error
	Close() error
}
