package app

import (
	"database/sql"
	authHander "dev/task-management/internal/modules/auth/handler"
	userRepo "dev/task-management/internal/modules/auth/repositories"
	authService "dev/task-management/internal/modules/auth/services"
	commentHandler "dev/task-management/internal/modules/comment/handler"
	commentRepo "dev/task-management/internal/modules/comment/repositories"
	commentService "dev/task-management/internal/modules/comment/services"
	healthHandler "dev/task-management/internal/modules/health/handler"
	taskHandler "dev/task-management/internal/modules/task/handler"
	taskRepo "dev/task-management/internal/modules/task/repositories"
	taskService "dev/task-management/internal/modules/task/services"
	workspaceHandler "dev/task-management/internal/modules/workspace/handler"
	workspaceRepo "dev/task-management/internal/modules/workspace/repositories"
	workspaceService "dev/task-management/internal/modules/workspace/services"
	"dev/task-management/internal/realtime"
	"dev/task-management/internal/router"
	"dev/task-management/pkg/cache"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

type App struct {
	Router *gin.Engine
}

func NewApp(database *sql.DB, redis *redis.Client) *App {

	hubConnectionWS := realtime.NewHubConnectionWS()

	redisService := cache.NewRedisCacheService(redis)

	// Subscriber: lắng nghe Redis Pub/Sub và forward vào Hub WS
	subscriber := realtime.NewNotificationSubscriber(hubConnectionWS, redisService)
	wsHandler := realtime.NewWSHandler(hubConnectionWS, subscriber)

	workspaceRepo := workspaceRepo.NewWorkspaceRepository(database)
	workspaceService := workspaceService.NewWorkspaceService(workspaceRepo)
	workspaceHandler := workspaceHandler.NewWorkspaceHandler(workspaceService)

	taskRepo := taskRepo.NewTaskRepository(database)
	taskService := taskService.NewTaskService(taskRepo, workspaceService, redisService)
	taskHandler := taskHandler.NewTaskHandler(taskService)

	userRepo := userRepo.NewUserRepository(database)
	authService := authService.NewAuthService(userRepo, workspaceService, redisService)
	authHander := authHander.NewAuthHandler(authService)

	commentRepo := commentRepo.NewCommentRepository(database)
	commentService := commentService.NewCommentService(commentRepo, redisService)
	commentHandler := commentHandler.NewCommentHandler(commentService)

	healthHandler := healthHandler.NewHealthHandler(redis, database)

	r := router.SetupRouter(router.RouterDependencies{
		TaskHandler:      taskHandler,
		AuthHander:       authHander,
		WorkspaceHandler: workspaceHandler,
		CommentHandler:   commentHandler,
		WSHandler:        wsHandler,
		HealthHandler:    healthHandler,
	})

	return &App{
		Router: r,
	}
}
