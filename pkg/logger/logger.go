package logger

import (
	"strings"
	"sync"

	"go.uber.org/zap"
)

var (
	mu       sync.RWMutex
	instance *zap.Logger
)

func Init(env string) {
	var (
		log *zap.Logger
		err error
	)

	switch strings.ToLower(env) {
	case "production":
		log, err = zap.NewProduction()
	default:
		cfg := zap.NewDevelopmentConfig()
		log, err = cfg.Build(zap.AddStacktrace(zap.ErrorLevel))
	}

	if err != nil {
		log = zap.NewNop()
	}

	mu.Lock()
	instance = log
	mu.Unlock()
}

func Info(message string, fields ...zap.Field) {
	GetLogger().Info(message, fields...)
}

func Error(message string, fields ...zap.Field) {
	GetLogger().Error(message, fields...)
}

func Warn(message string, fields ...zap.Field) {
	GetLogger().Warn(message, fields...)
}

func Debug(message string, fields ...zap.Field) {
	GetLogger().Debug(message, fields...)
}

func Sync() {
	_ = GetLogger().Sync()
}

func GetLogger() *zap.Logger {
	mu.RLock()
	log := instance
	mu.RUnlock()

	if log != nil {
		return log
	}

	mu.Lock()
	defer mu.Unlock()

	if instance == nil {
		instance = zap.NewNop()
	}

	return instance
}
