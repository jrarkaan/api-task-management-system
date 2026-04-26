package routes

import (
	"github.com/gin-gonic/gin"

	"api-task-management-system/modules/accounts/v1/deliveries"
)

func Init(router *gin.RouterGroup, authDelivery *deliveries.AuthDelivery) {
	auth := router.Group("/auth")
	auth.POST("/register", authDelivery.Register)
	auth.POST("/login", authDelivery.Login)
}
