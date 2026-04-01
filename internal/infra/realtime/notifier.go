package realtime

import (
	"context"
	"log/slog"

	"apiRest/internal/entity"
)

type Notifier struct {
	Hub    *Hub
	Logger *slog.Logger
}

func (n *Notifier) Send(ctx context.Context, notification entity.Notification) error {
	_ = ctx
	if n.Hub == nil {
		return nil
	}
	if err := n.Hub.Broadcast(notification.UserID, notification); err != nil {
		if n.Logger != nil {
			n.Logger.Error("websocket broadcast failed", "notification_id", notification.ID, "error", err)
		}
		return err
	}
	return nil
}
