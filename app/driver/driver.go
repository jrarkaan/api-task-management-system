package driver

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"

	"api-task-management-system/app/config"
	"api-task-management-system/app/middleware"
	"api-task-management-system/pkg/logger"
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

	router := gin.New()
	router.Use(middleware.Recovery())
	router.Use(middleware.RequestLogger())

	return &Driver{
		router: router,
		db:     db,
		cfg:    cfg,
	}
}

func (d *Driver) Run() error {
	address := fmt.Sprintf(":%s", d.cfg.AppPort)
	logger.Info("starting http server", zap.String("address", address))

	return d.router.Run(address)
}

func InitLogger(cfg config.Config) {
	logger.Init(cfg.AppEnv)
}
