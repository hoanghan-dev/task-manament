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
	database, err := config.ConnectPostgres()

	if err != nil {
		log.Fatalf("Connet database faild: %v \n", err)
	}

	defer database.Close()

	redis, err := config.NewRedisClient()

	if err != nil {
		log.Fatalf("Connet redis faild: %v \n", err)
	}

	defer redis.Close()

	notiRepo := repositories.NewNotificationRepository(database)
	notiService := services.NewNotificationService(notiRepo)

	ctx, canncel := context.WithCancel(context.Background())

	defer canncel()

	cache := cache.NewRedisCacheService(redis)

	workerNumber := 4

	for i := 1; i <= workerNumber; i++ {
		worker := worker.NewNotificationWorker(notiService, cache)
		go worker.StartCreateNotification(ctx)
	}

	// Giữ process sống, chờ Ctrl+C hoặc stop signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	<-quit

	log.Println("Shutting down notification worker...")
}
