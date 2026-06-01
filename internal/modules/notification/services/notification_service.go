package services

import (
	"context"
	"dev/task-management/internal/modules/notification/entities"
	"dev/task-management/internal/modules/notification/jobs"
	"dev/task-management/internal/modules/notification/repositories"
	"dev/task-management/internal/realtime"
	"dev/task-management/pkg/cache"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type NotificationService interface {
	CreateNotification(ctx context.Context, SenderId uuid.UUID, ReceiverId uuid.UUID, TaskId uuid.UUID) error
	SendAssignTaskNotification(ctx context.Context, userId uuid.UUID, assignTaskJob *jobs.AssignTaskJob) error
	SendUpdateStatusNotification(ctx context.Context, userId uuid.UUID, updateStatusJob *jobs.UpdateStatusJob) error
}

// --- Implementation cho API Server (dùng Hub WS trực tiếp) ---

type notificationService struct {
	repo        repositories.NotificationRepository
	redisClient *cache.RedisCacheService
}

func NewNotificationService(repo repositories.NotificationRepository, redisClient *cache.RedisCacheService) NotificationService {
	return &notificationService{
		repo:        repo,
		redisClient: redisClient,
	}
}

func (s *notificationService) SendAssignTaskNotification(ctx context.Context, userId uuid.UUID, assignTaskJob *jobs.AssignTaskJob) error {
	channel := fmt.Sprintf("notify:assign:%s", userId.String())
	payload := &realtime.Event{
		Type: "task.assigned",
		Data: assignTaskJob,
	}
	return s.redisClient.Publish(ctx, channel, payload)
}

func (s *notificationService) SendUpdateStatusNotification(ctx context.Context, userId uuid.UUID, updateStatusJob *jobs.UpdateStatusJob) error {
	channel := fmt.Sprintf("notify:status:%s", userId.String())
	payload := &realtime.Event{
		Type: "task.updated",
		Data: updateStatusJob,
	}
	return s.redisClient.Publish(ctx, channel, payload)
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
