package main

import (
	"dev/task-management/internal/app"
	"dev/task-management/internal/config"
	"fmt"
	"log"
)

func main() {
	fmt.Println("App running.....")
	database, err := config.ConnectPostgres()

	if err != nil {
		log.Fatalf("Connet database faild: %v \n", err)
	}

	defer database.Close()

	redis, err := config.NewRedisClient()

	if err != nil {
		log.Printf("Connet redis faild: %v \n", err)
	}

	defer redis.Close()

	appli := app.NewApp(database, redis)

	if err := appli.Router.Run(":8080"); err != nil {
		log.Fatalf("failed to start server: %v \n", err)
	}
}
