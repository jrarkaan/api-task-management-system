package driver

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"api-task-management-system/app/config"
)

type Driver struct {
	router *gin.Engine
	db     *gorm.DB
	cfg    config.Config
}

func New(db *gorm.DB, cfg config.Config) *Driver {
	if cfg.AppEnv == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	return &Driver{
		router: gin.Default(),
		db:     db,
		cfg:    cfg,
	}
}

func (d *Driver) Run() error {
	return d.router.Run(fmt.Sprintf(":%s", d.cfg.AppPort))
}
