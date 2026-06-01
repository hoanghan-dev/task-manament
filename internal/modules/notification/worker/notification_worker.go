package worker

import (
	"context"
	"dev/task-management/internal/modules/notification/jobs"
	"dev/task-management/internal/modules/notification/services"
	"dev/task-management/pkg/cache"
	"fmt"
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

func (w *NotificationWorker) StartCreateNotification(ctx context.Context) {
	var job jobs.NotificationJob
	queueKey := "queue:notifications"
	for {
		err := w.redisClient.Pop(ctx, queueKey, &job)

		if err != nil {
			fmt.Println("Pop notification job error:", err)
			continue
		}

		err = w.notiService.CreateNotification(ctx, job.SenderId, job.ReceiverId, job.TaskId)

		if err != nil {
			w.retryCreateNotiJob(ctx, job)
		}
		fmt.Println("Create notification successfully with receiver id: ", job.ReceiverId)
	}
}

func (w *NotificationWorker) retryCreateNotiJob(ctx context.Context, job jobs.NotificationJob) {

	var err error

	for attempt := 1; attempt <= 3; attempt++ {
		err = w.notiService.CreateNotification(ctx, job.SenderId, job.ReceiverId, job.TaskId)
		if err == nil {
			return
		}
		fmt.Println("Retry save notification, attempt:", attempt, "error:", err)
		if attempt < 3 {
			time.Sleep(300 * time.Millisecond)
		}
	}

	fmt.Printf("save notification failed after 3 attempts: %v\n", err)
}
