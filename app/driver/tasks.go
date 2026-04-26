package driver

import (
	"github.com/gin-gonic/gin"

	"api-task-management-system/app/middleware"
	taskDeliveries "api-task-management-system/modules/tasks/v1/deliveries"
	taskRepositories "api-task-management-system/modules/tasks/v1/repositories"
	taskRoutes "api-task-management-system/modules/tasks/v1/routes"
	taskUsecases "api-task-management-system/modules/tasks/v1/usecases"
)

func (d *Driver) registerTaskRoutes(api *gin.RouterGroup) {
	taskRepository := taskRepositories.NewTaskRepository(d.db, d.logger)
	taskUsecase := taskUsecases.NewTaskUsecase(taskRepository, d.txManager, d.logger)
	taskDelivery := taskDeliveries.NewTaskDelivery(taskUsecase)
	authMiddleware := middleware.Auth(d.cfg.JWTSecret)

	taskRoutes.Init(api, taskDelivery, authMiddleware)
}
