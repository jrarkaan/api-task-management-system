package middleware

import (
	"net/http"
	"runtime/debug"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"api-task-management-system/pkg/apiresponse"
	"api-task-management-system/pkg/logger"
)

func RequestLogger() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		start := time.Now()

		ctx.Next()

		latency := time.Since(start)
		status := ctx.Writer.Status()

		fields := []zap.Field{
			zap.String("method", ctx.Request.Method),
			zap.String("path", ctx.Request.URL.Path),
			zap.Int("status", status),
			zap.Duration("latency", latency),
			zap.String("client_ip", ctx.ClientIP()),
			zap.String("user_agent", ctx.Request.UserAgent()),
		}

		if status >= http.StatusInternalServerError {
			logger.Error("http request completed", fields...)
			return
		}

		if status >= http.StatusBadRequest {
			logger.Warn("http request completed", fields...)
			return
		}

		logger.Info("http request completed", fields...)
	}
}

func Recovery() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		defer func() {
			if recovered := recover(); recovered != nil {
				logger.Error(
					"panic recovered",
					zap.Any("error", recovered),
					zap.String("method", ctx.Request.Method),
					zap.String("path", ctx.Request.URL.Path),
					zap.String("client_ip", ctx.ClientIP()),
					zap.ByteString("stack", debug.Stack()),
				)

				apiresponse.ServerError(ctx, "internal server error")
				ctx.Abort()
			}
		}()

		ctx.Next()
	}
}
