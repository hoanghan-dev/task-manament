package entities

import (
	"time"

	"github.com/google/uuid"
)

type Notification struct {
	NotificationId uuid.UUID
	SenderId       uuid.UUID
	ReceiverId     uuid.UUID
	TaskId         uuid.UUID
	Message        string
	CreateAt       time.Time
}

func NewNotification(notificationId uuid.UUID, senderId uuid.UUID, receiverId uuid.UUID, taskId uuid.UUID, message string, createAt time.Time) *Notification {
	return &Notification{
		NotificationId: notificationId,
		SenderId:       senderId,
		ReceiverId:     receiverId,
		TaskId:         taskId,
		Message:        message,
		CreateAt:       createAt,
	}
}
