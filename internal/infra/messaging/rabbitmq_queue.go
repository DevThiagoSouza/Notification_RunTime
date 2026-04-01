package messaging

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"apiRest/internal/entity"
	amqp "github.com/rabbitmq/amqp091-go"
)

type RabbitMQQueue struct {
	conn      *amqp.Connection
	channel   *amqp.Channel
	queueName string
}

func NewRabbitMQQueue(url, queueName string) (*RabbitMQQueue, error) {
	if url == "" {
		return nil, errors.New("rabbitmq url is empty")
	}
	if queueName == "" {
		queueName = "notifications"
	}

	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, fmt.Errorf("rabbitmq dial: %w", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		_ = conn.Close()
		return nil, fmt.Errorf("rabbitmq channel: %w", err)
	}

	_, err = ch.QueueDeclare(queueName, true, false, false, false, nil)
	if err != nil {
		_ = ch.Close()
		_ = conn.Close()
		return nil, fmt.Errorf("rabbitmq declare queue: %w", err)
	}

	return &RabbitMQQueue{conn: conn, channel: ch, queueName: queueName}, nil
}

func (q *RabbitMQQueue) Publish(ctx context.Context, notification *entity.Notification) error {
	if notification == nil {
		return errors.New("notification is nil")
	}

	body, err := json.Marshal(notification)
	if err != nil {
		return err
	}

	return q.channel.PublishWithContext(ctx, "", q.queueName, false, false, amqp.Publishing{
		ContentType: "application/json",
		Body:        body,
	})
}

func (q *RabbitMQQueue) Consume(ctx context.Context, handler func(context.Context, entity.Notification) error) error {
	deliveries, err := q.channel.Consume(q.queueName, "", false, false, false, false, nil)
	if err != nil {
		return err
	}

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case d, ok := <-deliveries:
				if !ok {
					return
				}

				var n entity.Notification
				if err := json.Unmarshal(d.Body, &n); err != nil {
					_ = d.Nack(false, false)
					continue
				}

				if err := handler(ctx, n); err != nil {
					_ = d.Nack(false, true)
					continue
				}
				_ = d.Ack(false)
			}
		}
	}()

	return nil
}

func (q *RabbitMQQueue) Close() error {
	if q.channel != nil {
		_ = q.channel.Close()
	}
	if q.conn != nil {
		return q.conn.Close()
	}
	return nil
}
