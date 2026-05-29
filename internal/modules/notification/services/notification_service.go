package services

import (
	"context"
	"dev/task-management/internal/modules/notification/entities"
	"dev/task-management/internal/modules/notification/repositories"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type NotificationService interface {
	CreateNotification(ctx context.Context, SenderId uuid.UUID, ReceiverId uuid.UUID, TaskId uuid.UUID) error
	SendNotification(noti *entities.Notification)
}

type notificationService struct {
	repo repositories.NotificationRepository
}

func NewNotificationService(repo repositories.NotificationRepository) NotificationService {
	return &notificationService{
		repo: repo,
	}
}

func (s *notificationService) SendNotification(noti *entities.Notification) {
	fmt.Printf("[INFO] Notification has been sent with message: %v", noti.Message)
}

func (s *notificationService) CreateNotification(ctx context.Context, SenderId uuid.UUID, ReceiverId uuid.UUID, TaskId uuid.UUID) error {
	message := "You have been assigned a new task"

	noti := entities.NewNotification(uuid.New(), SenderId, ReceiverId, TaskId, message, time.Now())

	err := s.repo.SaveNotification(ctx, noti)
	if err != nil {
		return err
	}

	return nil

}
