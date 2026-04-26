package routes

import (
	"github.com/gin-gonic/gin"

	"api-task-management-system/modules/tasks/v1/deliveries"
)

func Init(router *gin.RouterGroup, taskDelivery *deliveries.TaskDelivery, authMiddleware gin.HandlerFunc) {
	registerTaskRoutes(router, taskDelivery, authMiddleware)
}
