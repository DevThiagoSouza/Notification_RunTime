package entity

import "errors"

var ErrNotificationNotFound = errors.New("Invalid notification: Required fields missing.")

const (
	StatusPending    = "pending"
	StatusProcessing = "processing"
	StatusSent       = "sent"
	StatusFailed     = "failed"
)

type Notification struct {
	ID      string
	UserID  string
	Content string
	Type    string
	Status  string
}

func (n *Notification) Validate() error {
	if n.ID == "" || n.UserID == "" || n.Content == "" || n.Type == "" {
		return ErrNotificationNotFound
	}

	return nil
}
