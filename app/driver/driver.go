package driver

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"

	"api-task-management-system/app/config"
	"api-task-management-system/app/middleware"
	dbpkg "api-task-management-system/pkg/db"
	loggerpkg "api-task-management-system/pkg/logger"
)

type Driver struct {
	router    *gin.Engine
	db        *gorm.DB
	cfg       config.Config
	logger    *zap.Logger
	txManager *dbpkg.TransactionManager
}

func New(db *gorm.DB, cfg config.Config) *Driver {
	if cfg.AppEnv == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.New()
	router.Use(middleware.Recovery())
	router.Use(middleware.RequestLogger())

	log := loggerpkg.GetLogger()
	txManager := dbpkg.NewTransactionManager(db)

	return &Driver{
		router:    router,
		db:        db,
		cfg:       cfg,
		logger:    log,
		txManager: txManager,
	}
}

func (d *Driver) Run() error {
	address := fmt.Sprintf(":%s", d.cfg.AppPort)
	loggerpkg.Info("starting http server", zap.String("address", address))

	return d.router.Run(address)
}

func InitLogger(cfg config.Config) {
	loggerpkg.Init(cfg.AppEnv)
}
