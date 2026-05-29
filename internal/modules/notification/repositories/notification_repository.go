package repositories

import (
	"context"
	"database/sql"
	"dev/task-management/internal/modules/notification/entities"
	"errors"

	"github.com/google/uuid"
)

type NotificationRepository interface {
	FindReceivedNotifications(ctx context.Context, receiverId uuid.UUID) ([]*entities.Notification, error)
	SaveNotification(ctx context.Context, noti *entities.Notification) error
}

type notificationRepository struct {
	db *sql.DB
}

func NewNotificationRepository(db *sql.DB) NotificationRepository {
	return &notificationRepository{
		db: db,
	}
}

func (r *notificationRepository) FindReceivedNotifications(ctx context.Context, receiverId uuid.UUID) ([]*entities.Notification, error) {
	sqlQuery := `select notification_id, sender_id, receiver_id, task_id, message, create_at 
				from notifications where receiver_id = $1`
	rows, err := r.db.QueryContext(ctx, sqlQuery, receiverId)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	notifications := make([]*entities.Notification, 0)

	for rows.Next() {
		var notification entities.Notification
		err := rows.Scan(
			&notification.NotificationId,
			&notification.SenderId,
			&notification.ReceiverId,
			&notification.TaskId,
			&notification.Message,
			&notification.CreateAt)
		if err != nil {
			return nil, err
		}

		notifications = append(notifications, &notification)

	}
	return notifications, nil
}

func (r *notificationRepository) SaveNotification(ctx context.Context, noti *entities.Notification) error {
	sqlQuery := `insert into notifications (notification_id, sender_id, receiver_id, task_id, message, create_at)
				values ($1, $2, $3, $4, $5, $6)`

	result, err := r.db.ExecContext(ctx, sqlQuery, noti.NotificationId, noti.SenderId, noti.ReceiverId, noti.TaskId, noti.Message, noti.CreateAt)
	if err != nil {
		return err
	}

	numOfRow, err := result.RowsAffected()

	if err != nil {
		return err
	}

	if numOfRow == 0 {
		return errors.New("faild to save notification")
	}

	return nil
}
