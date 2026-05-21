package router

import (
	"dev/task-management/internal/middleware"
	authHandler "dev/task-management/internal/modules/auth/handler"
	taskHandler "dev/task-management/internal/modules/task/handler"

	"github.com/gin-gonic/gin"
)

type RouterDependencies struct {
	TaskHandler *taskHandler.TaskHandler
	AuthHander  *authHandler.AuthHander
}

func SetupRouter(deps RouterDependencies) *gin.Engine {
	r := gin.New()

	r.Use(middleware.LoggingMiddleware,
		middleware.RequestIdMiddleware,
		middleware.AuthMiddleware,
		gin.Recovery())

	api := r.Group("api/")
	SetupTaskRouter(api, deps.TaskHandler)
	SetupAuthRouter(api, deps.AuthHander)
	return r
}

func SetupTaskRouter(api *gin.RouterGroup, taskHandler *taskHandler.TaskHandler) {
	tasks := api.Group("tasks/")
	// URL mapping với api tasks/
	{
		tasks.GET("/", taskHandler.GetAllTask)
		tasks.GET("/:id", taskHandler.GetTask)
		tasks.POST("/", taskHandler.CreateTask)
		tasks.PUT("/:id", taskHandler.UpdateTask)
		tasks.DELETE("/:id", taskHandler.DeleteTask)
	}
}

func SetupAuthRouter(api *gin.RouterGroup, authHander *authHandler.AuthHander) {
	auth := api.Group("auth/")
	{
		auth.POST("/register", authHander.Register)
		auth.POST("/login", authHander.Login)
	}
}
