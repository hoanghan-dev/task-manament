package main

import (
	"dev/task-management/internal/router"
	"fmt"
)

func main() {
	fmt.Println("App running.....")
	r := router.SetupRouter()
	r.Run(":8080")
}
