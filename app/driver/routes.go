package driver

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (d *Driver) RegisterRoutes() {
	d.router.GET("/health", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{"message": "ok"})
	})

	api := d.router.Group("/api")

	d.registerAccountRoutes(api)
	d.registerTaskRoutes(api)
}
