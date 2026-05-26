package router

import (
	"dev/task-management/internal/middleware"
	authHandler "dev/task-management/internal/modules/auth/handler"
	taskHandler "dev/task-management/internal/modules/task/handler"
	workspaceHandler "dev/task-management/internal/modules/workspace/handler"

	"github.com/gin-gonic/gin"
)

type RouterDependencies struct {
	TaskHandler      *taskHandler.TaskHandler
	AuthHander       *authHandler.AuthHander
	WorkspaceHandler *workspaceHandler.WorkspaceHandler
}

func SetupRouter(deps RouterDependencies) *gin.Engine {
	r := gin.New()

	r.Use(middleware.LoggingMiddleware,
		middleware.RequestIdMiddleware,
		gin.Recovery())

	api := r.Group("api/")
	SetupAuthRouter(api, deps.AuthHander)

	protected := api.Group("", middleware.AuthMiddleware)
	SetupTaskRouter(protected, deps.TaskHandler)
	SetupWorkspaceRouter(protected, deps.WorkspaceHandler)
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

func SetupWorkspaceRouter(api *gin.RouterGroup, workspaceHandler *workspaceHandler.WorkspaceHandler) {
	workspace := api.Group("workspaces/")
	{
		workspace.GET("/", workspaceHandler.GetWorkspace)
		workspace.PUT("/", workspaceHandler.UpdateWorkspace)
		workspace.DELETE("/:wpId", workspaceHandler.DeleteWorkspace)
	}
}
