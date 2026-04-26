package driver

import (
	"net/http"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func (d *Driver) RegisterRoutes() {
	d.router.GET("/health", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{"message": "ok"})
	})

	d.router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	apiAccount := d.router.Group("/accounts")
	apiTasks := d.router.Group("")

	d.registerAccountRoutes(apiAccount)
	d.registerTaskRoutes(apiTasks)
}
