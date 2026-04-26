package main

import (
	"go.uber.org/zap"

	"api-task-management-system/app/config"
	"api-task-management-system/app/driver"
	_ "api-task-management-system/docs"
	db_pg "api-task-management-system/pkg/db/pg"
	"api-task-management-system/pkg/logger"
)

// @title Task Management System API
// @version 1.0
// @description REST API for user authentication and personal task management.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.email support@example.com

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host localhost:8080
// @BasePath /
// @schemes http

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.
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
