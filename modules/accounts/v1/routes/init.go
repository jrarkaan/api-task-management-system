package routes

import (
	"github.com/gin-gonic/gin"

	"api-task-management-system/modules/accounts/v1/deliveries"
)

func Init(router *gin.RouterGroup, authDelivery *deliveries.AuthDelivery) {
	registerAuthRoutes(router, authDelivery)
}
