package main

import (
	"log"

	"api-task-management-system/app/config"
	"api-task-management-system/app/driver"
	"api-task-management-system/pkg/db"
)

func main() {
	cfg := config.Load()

	database, err := db.NewPostgresConnection(cfg)
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}

	app := driver.New(database, cfg)
	app.RegisterRoutes()

	if err := app.Run(); err != nil {
		log.Fatalf("failed to run server: %v", err)
	}
}
