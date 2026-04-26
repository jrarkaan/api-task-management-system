package main

import (
	"go.uber.org/zap"

	"api-task-management-system/app/config"
	"api-task-management-system/app/driver"
	db_pg "api-task-management-system/pkg/db/pg"
	"api-task-management-system/pkg/logger"
)

func main() {
	cfg := config.Load()
	driver.InitLogger(cfg)
	defer logger.Sync()

	database, err := db_pg.InitDB(&cfg)
	if err != nil {
		logger.Error("failed to connect database", zap.Error(err))
		return
	}

	app := driver.New(database, cfg)
	app.RegisterRoutes()

	if err := app.Run(); err != nil {
		logger.Error("failed to run server", zap.Error(err))
	}
}
