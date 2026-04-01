package repository

import (
	"apiRest/internal/entity"
	"context"
	"database/sql"
	"errors"
)

type PostgresNotificationRepository struct {
	DB *sql.DB
}

func NewPostgresNotificationRepository(db *sql.DB) *PostgresNotificationRepository {
	return &PostgresNotificationRepository{DB: db}
}

func (r *PostgresNotificationRepository) Save(ctx context.Context, notification *entity.Notification) error {
	query := "INSERT INTO notifications (id, user_id, content, type, status) VALUES ($1, $2, $3, $4, $5)"
	_, err := r.DB.ExecContext(ctx, query, notification.ID, notification.UserID, notification.Content, notification.Type, notification.Status)
	return err
}

func (r *PostgresNotificationRepository) UpdateStatus(ctx context.Context, notificationID, status string) error {
	query := "UPDATE notifications SET status = $1 WHERE id = $2"
	result, err := r.DB.ExecContext(ctx, query, status, notificationID)
	if err != nil {
		return err
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if affected == 0 {
		return errors.New("notification not found")
	}

	return nil
}
