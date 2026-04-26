package driver

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"api-task-management-system/app/middleware"
	accountDeliveries "api-task-management-system/modules/accounts/v1/deliveries"
	accountRepositories "api-task-management-system/modules/accounts/v1/repositories"
	accountRoutes "api-task-management-system/modules/accounts/v1/routes"
	accountUsecases "api-task-management-system/modules/accounts/v1/usecases"
	taskDeliveries "api-task-management-system/modules/tasks/v1/deliveries"
	taskRepositories "api-task-management-system/modules/tasks/v1/repositories"
	taskRoutes "api-task-management-system/modules/tasks/v1/routes"
	taskUsecases "api-task-management-system/modules/tasks/v1/usecases"
)

func (d *Driver) RegisterRoutes() {
	d.router.GET("/health", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{"message": "ok"})
	})

	api := d.router.Group("/api")

	userRepository := accountRepositories.NewUserRepository(d.db)
	authUsecase := accountUsecases.NewAuthUsecase(userRepository, d.cfg.JWTSecret, d.cfg.JWTExpiresHours)
	authDelivery := accountDeliveries.NewAuthDelivery(authUsecase)
	accountRoutes.Init(api, authDelivery)

	taskRepository := taskRepositories.NewTaskRepository(d.db)
	taskUsecase := taskUsecases.NewTaskUsecase(taskRepository)
	taskDelivery := taskDeliveries.NewTaskDelivery(taskUsecase)
	authMiddleware := middleware.Auth(d.cfg.JWTSecret)
	taskRoutes.Init(api, taskDelivery, authMiddleware)
}
