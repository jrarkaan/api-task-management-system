package driver

import (
	"net/http"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"api-task-management-system/pkg/apiresponse"
)

func (d *Driver) RegisterRoutes() {
	d.router.NoRoute(func(ctx *gin.Context) {
		if ctx.Request.Method == http.MethodOptions {
			ctx.Status(http.StatusNoContent)
			return
		}
		apiresponse.NotFound(ctx, nil, "route not found")
	})

	d.router.GET("/health", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{"message": "ok"})
	})

	d.router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	apiAccount := d.router.Group("/accounts")
	apiTasks := d.router.Group("")

	d.registerAccountRoutes(apiAccount)
	d.registerTaskRoutes(apiTasks)
}
