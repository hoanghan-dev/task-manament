package worker

import (
	"context"
	"dev/task-management/internal/modules/notification/jobs"
	"dev/task-management/internal/modules/notification/services"
	"dev/task-management/pkg/cache"
	"fmt"
	"log"
	"time"
)

type NotificationWorker struct {
	notiService services.NotificationService
	redisClient *cache.RedisCacheService
}

func NewNotificationWorker(notiService services.NotificationService, redisClient *cache.RedisCacheService) *NotificationWorker {
	return &NotificationWorker{
		notiService: notiService,
		redisClient: redisClient,
	}
}

func (w *NotificationWorker) StartAssignTaskNotification(ctx context.Context) {
	queueKey := "queue:notifications:assign"
	for {
		var job jobs.AssignTaskJob
		err := w.redisClient.Pop(ctx, queueKey, &job)

		if err != nil {
			fmt.Println("Pop notification job error:", err)
			continue
		}
		w.handleAssignTaskNotification(ctx, job)
	}
}

func (w *NotificationWorker) StartUpdateStatusNotification(ctx context.Context) {
	queueKey := "queue:notifications:status"
	for {
		var job jobs.UpdateStatusJob
		err := w.redisClient.Pop(ctx, queueKey, &job)

		if err != nil {
			fmt.Println("Pop notification job error:", err)
			continue
		}
		w.handleUpdateStatusNotification(ctx, job)
	}
}

func (w *NotificationWorker) handleUpdateStatusNotification(ctx context.Context, job jobs.UpdateStatusJob) {
	err := w.notiService.SendUpdateStatusNotification(ctx, job.ReceiverId, &job)
	if err == nil {
		log.Println("Update status notification successfully with receiver id: ", job.ReceiverId)
		return
	}
	w.retryUpdateStatusNotiJob(ctx, job)
}

func (w *NotificationWorker) handleAssignTaskNotification(ctx context.Context, job jobs.AssignTaskJob) {
	w.StartCreateNotification(ctx, job)
	w.StartSendNotification(ctx, job)
}

func (w *NotificationWorker) StartCreateNotification(ctx context.Context, job jobs.AssignTaskJob) {
	err := w.notiService.CreateNotification(ctx, job.SenderId, job.ReceiverId, job.TaskId)

	if err == nil {
		log.Println("Create notification successfully with receiver id: ", job.ReceiverId)
		return
	}

	w.retryCreateNotiJob(ctx, job)
}

func (w *NotificationWorker) StartSendNotification(ctx context.Context, job jobs.AssignTaskJob) {
	err := w.notiService.SendAssignTaskNotification(ctx, job.ReceiverId, &job)

	if err == nil {
		log.Println("Send notification successfully with receiver id: ", job.ReceiverId)
		return
	}

	w.retrySendNotiJob(ctx, job)
}

func (w *NotificationWorker) retryCreateNotiJob(ctx context.Context, job jobs.AssignTaskJob) {

	var err error

	for attempt := 1; attempt <= 3; attempt++ {
		err = w.notiService.CreateNotification(ctx, job.SenderId, job.ReceiverId, job.TaskId)
		if err == nil {
			return
		}
		log.Println("Retry save notification, attempt:", attempt, "error:", err)
		if attempt < 3 {
			time.Sleep(1 * time.Second)
		}
	}
	log.Printf("Save notification failed after 3 attempts: %v\n", err)
}

func (w *NotificationWorker) retryUpdateStatusNotiJob(ctx context.Context, job jobs.UpdateStatusJob) {
	var err error
	for attempt := 1; attempt <= 3; attempt++ {
		err = w.notiService.SendUpdateStatusNotification(ctx, job.ReceiverId, &job)
		if err == nil {
			return
		}
		log.Println("Retry update status notification, attempt:", attempt, "error:", err)
		if attempt < 3 {
			time.Sleep(1 * time.Second)
		}
	}
	log.Printf("Update status notification failed after 3 attempts: %v\n", err)
}

func (w *NotificationWorker) retrySendNotiJob(ctx context.Context, job jobs.AssignTaskJob) {

	var err error

	for attempt := 1; attempt <= 3; attempt++ {
		err = w.notiService.SendAssignTaskNotification(ctx, job.ReceiverId, &job)
		if err == nil {
			return
		}
		log.Println("Retry send notification, attempt:", attempt, "error:", err)
		if attempt < 3 {
			time.Sleep(1 * time.Second)
		}
	}
	log.Printf("Send notification failed after 3 attempts: %v\n", err)
}
