package router

import (
	"dev/task-management/internal/handler"

	"github.com/gin-gonic/gin"
)

func SetupRouter(taskHandler handler.TaskHandler) *gin.Engine {
	router := gin.Default() // khởi tạo router

	api := router.Group("api/") // khởi tạo prefix api

	// GET mapping với api tasks/
	api.GET("/tasks", taskHandler.GetAllTask)
	api.GET("/tasks/:id", taskHandler.GetTask)
	api.POST("/tasks", taskHandler.CreateTask)
	api.PUT("/tasks/:id", taskHandler.UpdateTask)
	api.DELETE("/tasks/:id", taskHandler.DeleteTask)

	return router
}
