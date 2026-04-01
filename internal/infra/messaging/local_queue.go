package messaging

import (
	"context"
	"errors"

	"apiRest/internal/entity"
)

type LocalQueue struct {
	ch chan entity.Notification
}

func NewLocalQueue(buffer int) *LocalQueue {
	if buffer <= 0 {
		buffer = 100
	}
	return &LocalQueue{ch: make(chan entity.Notification, buffer)}
}

func (q *LocalQueue) Publish(ctx context.Context, notification *entity.Notification) error {
	if notification == nil {
		return errors.New("notification is nil")
	}

	select {
	case <-ctx.Done():
		return ctx.Err()
	case q.ch <- *notification:
		return nil
	}
}

func (q *LocalQueue) Consume(ctx context.Context, handler func(context.Context, entity.Notification) error) error {
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case n := <-q.ch:
				_ = handler(ctx, n)
			}
		}
	}()
	return nil
}

func (q *LocalQueue) Close() error {
	return nil
}
