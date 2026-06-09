package router

import (
	"dev/task-management/internal/middleware"
	authHandler "dev/task-management/internal/modules/auth/handler"
	commentHandler "dev/task-management/internal/modules/comment/handler"
	healthHandler "dev/task-management/internal/modules/health/handler"
	taskHandler "dev/task-management/internal/modules/task/handler"
	workspaceHandler "dev/task-management/internal/modules/workspace/handler"
	"dev/task-management/internal/realtime"

	"github.com/gin-gonic/gin"
)

type RouterDependencies struct {
	TaskHandler      *taskHandler.TaskHandler
	AuthHander       *authHandler.AuthHander
	WorkspaceHandler *workspaceHandler.WorkspaceHandler
	CommentHandler   *commentHandler.CommentHandler
	WSHandler        *realtime.WSHandler
	HealthHandler    *healthHandler.HealthHandler
}

func SetupRouter(deps RouterDependencies) *gin.Engine {
	r := gin.New()

	r.Use(middleware.LoggingMiddleware,
		middleware.RequestIdMiddleware,
		gin.Recovery())

	api := r.Group("api/")
	SetupAuthRouter(api, deps.AuthHander)
	SetupHealthRouter(api, deps.HealthHandler)

	protected := api.Group("", middleware.AuthMiddleware)
	SetupTaskRouter(protected, deps.TaskHandler)
	SetupWorkspaceRouter(protected, deps.WorkspaceHandler)
	SetupCommentRouter(protected, deps.CommentHandler)
	SetupWSRouter(protected, deps.WSHandler)
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
		tasks.PATCH("/assign", taskHandler.AssignTask)
		tasks.PATCH("/:id/status", taskHandler.UpdateTaskStatus)
	}
}

func SetupAuthRouter(api *gin.RouterGroup, authHander *authHandler.AuthHander) {
	auth := api.Group("auth/")
	{
		auth.POST("/register", authHander.Register)
		auth.POST("/login", authHander.Login)
		auth.POST("/refresh", authHander.RefreshToken)
	}
	auth.Group("", middleware.AuthMiddleware).GET("/logout", authHander.Logout)

}

func SetupWorkspaceRouter(api *gin.RouterGroup, workspaceHandler *workspaceHandler.WorkspaceHandler) {
	workspace := api.Group("workspaces/")
	{
		workspace.GET("/", workspaceHandler.GetWorkspace)
		workspace.PUT("/", workspaceHandler.UpdateWorkspace)
		workspace.DELETE("/:wpId", workspaceHandler.DeleteWorkspace)
	}
}

func SetupCommentRouter(api *gin.RouterGroup, handler *commentHandler.CommentHandler) {
	tasks := api.Group("tasks/")
	{
		tasks.POST("/:id/comments", handler.CreateComment)
		tasks.GET("/:id/comments", handler.GetComments)
	}

	comments := api.Group("comments/")
	{
		comments.DELETE("/:commentId", handler.DeleteComment)
	}
}

func SetupWSRouter(api *gin.RouterGroup, wsHandler *realtime.WSHandler) {
	ws := api.Group("ws/")
	{
		ws.GET("/", wsHandler.Connect)
	}
}

func SetupHealthRouter(api *gin.RouterGroup, HealthHandler *healthHandler.HealthHandler) {
	health := api.Group("health/")
	{
		health.GET("/", HealthHandler.HealthCheck)
	}
}
