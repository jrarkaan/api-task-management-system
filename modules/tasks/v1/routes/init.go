package routes

import (
	"github.com/gin-gonic/gin"

	"api-task-management-system/modules/tasks/v1/deliveries"
)

func Init(router *gin.RouterGroup, taskDelivery *deliveries.TaskDelivery, authMiddleware gin.HandlerFunc) {
	tasks := router.Group("/tasks")
	tasks.Use(authMiddleware)
	tasks.GET("", taskDelivery.List)
	tasks.POST("", taskDelivery.Create)
	tasks.PUT("/:id", taskDelivery.Update)
	tasks.DELETE("/:id", taskDelivery.Delete)
}
