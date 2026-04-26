package driver

import (
	"github.com/gin-gonic/gin"

	accountDeliveries "api-task-management-system/modules/accounts/v1/deliveries"
	accountRepositories "api-task-management-system/modules/accounts/v1/repositories"
	accountRoutes "api-task-management-system/modules/accounts/v1/routes"
	accountUsecases "api-task-management-system/modules/accounts/v1/usecases"
)

func (d *Driver) registerAccountRoutes(api *gin.RouterGroup) {
	userRepository := accountRepositories.NewUserRepository(d.db, d.logger)
	authUsecase := accountUsecases.NewAuthUsecase(userRepository, d.txManager, d.logger, d.cfg.JWTSecret, d.cfg.JWTExpiresHours)
	authDelivery := accountDeliveries.NewAuthDelivery(authUsecase)

	accountRoutes.Init(api, authDelivery)
}
