package app

import (
	"database/sql"
	authHander "dev/task-management/internal/modules/auth/handler"
	userRepo "dev/task-management/internal/modules/auth/repositories"
	authService "dev/task-management/internal/modules/auth/services"
	notiRepo "dev/task-management/internal/modules/notification/repositories"
	notiService "dev/task-management/internal/modules/notification/services"
	taskHandler "dev/task-management/internal/modules/task/handler"
	taskRepo "dev/task-management/internal/modules/task/repositories"
	taskService "dev/task-management/internal/modules/task/services"
	workspaceHandler "dev/task-management/internal/modules/workspace/handler"
	workspaceRepo "dev/task-management/internal/modules/workspace/repositories"
	workspaceService "dev/task-management/internal/modules/workspace/services"
	"dev/task-management/internal/router"
	"dev/task-management/pkg/cache"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

type App struct {
	Router *gin.Engine
}

func NewApp(database *sql.DB, redis *redis.Client) *App {

	redisService := cache.NewRedisCacheService(redis)

	workspaceRepo := workspaceRepo.NewWorkspaceRepository(database)
	workspaceService := workspaceService.NewWorkspaceService(workspaceRepo)
	workspaceHandler := workspaceHandler.NewWorkspaceHandler(workspaceService)

	notiRepo := notiRepo.NewNotificationRepository(database)
	notiService := notiService.NewNotificationService(notiRepo)

	taskRepo := taskRepo.NewTaskRepository(database)
	taskService := taskService.NewTaskService(taskRepo, workspaceService, redisService, notiService)
	taskHandler := taskHandler.NewTaskHandler(taskService)

	userRepo := userRepo.NewUserRepository(database)
	authService := authService.NewAuthService(userRepo, workspaceService)
	authHander := authHander.NewAuthHandler(authService)
	r := router.SetupRouter(router.RouterDependencies{
		TaskHandler:      taskHandler,
		AuthHander:       authHander,
		WorkspaceHandler: workspaceHandler,
	})

	return &App{
		Router: r,
	}
}
