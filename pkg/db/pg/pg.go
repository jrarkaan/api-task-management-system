package db_pg

import (
	"fmt"
	"sync"
	"time"

	"api-task-management-system/app/config"
	"api-task-management-system/pkg/logger"
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
)

var (
	once        sync.Once
	instance    *gorm.DB
	instanceErr error
)

const (
	maxOpenConns    = 60
	connMaxLifetime = 120
	maxIdleConns    = 30
	connMaxIdleTime = 20
)

func InitDB(conf *config.Config) (*gorm.DB, error) {
	once.Do(func() {
		dsn := fmt.Sprintf(
			"host=%s port=%s user=%s dbname=%s password=%s sslmode=%s TimeZone=%s",
			conf.DBHost,
			conf.DBPort,
			conf.DBUser,
			conf.DBName,
			conf.DBPassword,
			conf.DBSSLMode,
			conf.DBTimezone,
		)

		db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
			Logger: gormlogger.Default.LogMode(gormlogger.Info),
		})
		if err != nil {
			instanceErr = fmt.Errorf("gorm.Open: %w", err)
			logger.Error("failed to open database connection", zap.Error(instanceErr))
			return
		}

		sqlDB, err := db.DB()
		if err != nil {
			instanceErr = fmt.Errorf("db.DB: %w", err)
			logger.Error("failed to get database sql instance", zap.Error(instanceErr))
			return
		}

		sqlDB.SetMaxOpenConns(maxOpenConns)
		sqlDB.SetConnMaxLifetime(connMaxLifetime * time.Second)
		sqlDB.SetMaxIdleConns(maxIdleConns)
		sqlDB.SetConnMaxIdleTime(connMaxIdleTime * time.Second)

		if err := sqlDB.Ping(); err != nil {
			instanceErr = fmt.Errorf("db.Ping: %w", err)
			logger.Error("failed to ping database", zap.Error(instanceErr))
			return
		}

		instance = db
		logger.Info(
			"database connection initialized",
			zap.String("host", conf.DBHost),
			zap.String("port", conf.DBPort),
			zap.String("database", conf.DBName),
			zap.Int("max_open_conns", maxOpenConns),
			zap.Int("max_idle_conns", maxIdleConns),
		)
	})

	if instanceErr != nil {
		return nil, instanceErr
	}

	if instance == nil {
		return nil, fmt.Errorf("database instance is nil")
	}

	return instance, nil
}

func GetDB() (*gorm.DB, error) {
	if instance == nil {
		return nil, fmt.Errorf("database is not initialized")
	}

	return instance, nil
}
