package app

import (
	"database/sql"
	authHander "dev/task-management/internal/modules/auth/handler"
	userRepo "dev/task-management/internal/modules/auth/repositories"
	authService "dev/task-management/internal/modules/auth/services"
	taskHandler "dev/task-management/internal/modules/task/handler"
	taskRepo "dev/task-management/internal/modules/task/repositories"
	taskService "dev/task-management/internal/modules/task/services"
	"dev/task-management/internal/router"

	"github.com/gin-gonic/gin"
)

type App struct {
	Router *gin.Engine
}

func NewApp(database *sql.DB) *App {
	taskRepo := taskRepo.NewTaskRepository(database)
	taskService := taskService.NewTaskService(taskRepo)
	taskHandler := taskHandler.NewTaskHandler(taskService)

	userRepo := userRepo.NewUserRepository(database)
	authService := authService.NewAuthService(userRepo)
	authHander := authHander.NewAuthHandler(authService)

	r := router.SetupRouter(router.RouterDependencies{
		TaskHandler: taskHandler,
		AuthHander:  authHander,
	})

	return &App{
		Router: r,
	}
}
