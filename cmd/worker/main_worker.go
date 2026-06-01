package main

import (
	"context"
	"dev/task-management/internal/config"
	"dev/task-management/internal/modules/notification/repositories"
	"dev/task-management/internal/modules/notification/services"
	"dev/task-management/internal/modules/notification/worker"
	"dev/task-management/pkg/cache"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	fmt.Println("Worker running.....")

	// Kết nối PostgreSQL để lưu thông báo vào database
	database, err := config.ConnectPostgres()
	if err != nil {
		log.Fatalf("Connet database faild: %v \n", err)
	}
	defer database.Close()

	// Kết nối Redis để đọc job từ queue và publish thông báo qua Pub/Sub
	redis, err := config.NewRedisClient()
	if err != nil {
		log.Fatalf("Connet redis faild: %v \n", err)
	}
	defer redis.Close()

	// Khởi tạo Redis cache service (dùng chung cho queue và pub/sub)
	redisService := cache.NewRedisCacheService(redis)

	// Khởi tạo notification repository và service
	// Service sẽ PUBLISH event lên Redis Pub/Sub thay vì gửi WebSocket trực tiếp
	// (vì Worker là process riêng, không có Hub WebSocket)
	notiRepo := repositories.NewNotificationRepository(database)
	notiService := services.NewNotificationService(notiRepo, redisService)

	// Tạo context để kiểm soát vòng đời của các worker goroutine
	// Khi cancel() được gọi, tất cả worker sẽ dừng lại
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Worker 1: Xử lý job "assign task" từ queue:notifications:assign
	// Lưu thông báo vào DB và publish event lên Redis channel notify:assign:{userId}
	go func() {
		assignTaskWorker := worker.NewNotificationWorker(notiService, redisService)
		assignTaskWorker.StartAssignTaskNotification(ctx)
	}()

	// Worker 2: Xử lý job "update status" từ queue:notifications:status
	// Publish event lên Redis channel notify:status:{userId}
	go func() {
		updateStatusWorker := worker.NewNotificationWorker(notiService, redisService)
		updateStatusWorker.StartUpdateStatusNotification(ctx)
	}()

	// Giữ process sống, chờ tín hiệu dừng từ hệ thống (Ctrl+C hoặc kill)
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down notification worker...")
}
