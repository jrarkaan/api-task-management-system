package driver

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (d *Driver) RegisterRoutes() {
	d.router.GET("/health", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{"message": "ok"})
	})

	apiAccount := d.router.Group("/accounts")
	apiTasks := d.router.Group("/tasks")

	d.registerAccountRoutes(apiAccount)
	d.registerTaskRoutes(apiTasks)
}
