package main

import (
	"dev/task-management/internal/handler"
	"dev/task-management/internal/repositories"
	"dev/task-management/internal/router"
	"dev/task-management/internal/services"
	"fmt"
)

func main() {
	fmt.Println("App running.....")
	repo := repositories.NewTaskRepository()
	service := services.NewTaskService(repo)
	handler := handler.NewTaskHandler(service)
	r := router.SetupRouter(*handler)
	r.Run(":8080")
}
