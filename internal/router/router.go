package router

import (
	"dev/task-management/internal/handler"
	"dev/task-management/internal/middleware"
	"dev/task-management/internal/repositories"
	"dev/task-management/internal/services"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	router := gin.New()         // khởi tạo router
	api := router.Group("api/") // khởi tạo prefix api
	api.Use(
		middleware.RequestIdMiddleware,
		middleware.LoggingMiddleware,
		gin.Recovery(),
	)

	SetupTaskRouter(api)
	SetupAssessmentRouter(api)
	return router
}

func SetupTaskRouter(api *gin.RouterGroup) {
	repo := repositories.NewTaskRepository()
	service := services.NewTaskService(repo)
	taskHandler := handler.NewTaskHandler(service)
	tasks := api.Group("tasks/")
	// URL mapping với api tasks/
	tasks.GET("/", taskHandler.GetAllTask)
	tasks.GET("/:id", taskHandler.GetTask)
	tasks.POST("/", taskHandler.CreateTask)
	tasks.PUT("/:id", taskHandler.UpdateTask)
	tasks.DELETE("/:id", taskHandler.DeleteTask)
}

func SetupAssessmentRouter(api *gin.RouterGroup) {
	repo := repositories.NewAssessmentRepository()
	services := services.NewAssessmentService(repo)
	handler := handler.NewAssessmentHandler(services)

	assessments := api.Group("assessments/")
	assessments.GET("/", handler.GetAssessmentList)
	assessments.GET("/:id", handler.GetAssessment)
	assessments.POST("/", handler.CreateAssessment)
	assessments.PUT("/:id", handler.UpdateAssessment)
	assessments.DELETE("/:id", handler.DeleteAssessment)
}
