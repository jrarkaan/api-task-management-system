package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"api-task-management-system/pkg/apiresponse"
	"api-task-management-system/pkg/helpers"
	"api-task-management-system/pkg/logger"
)

const UserIDContextKey = "user_id"

func Auth(jwtSecret string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authHeader := ctx.GetHeader("Authorization")
		if authHeader == "" {
			logUnauthorized(ctx, "authorization header is required")
			apiresponse.UnAuthorized(ctx, nil, "authorization header is required")
			ctx.Abort()
			return
		}

		parts := strings.Fields(authHeader)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
			logUnauthorized(ctx, "authorization header must use Bearer token")
			apiresponse.UnAuthorized(ctx, nil, "authorization header must use Bearer token")
			ctx.Abort()
			return
		}

		claims, err := helpers.ParseJWT(parts[1], jwtSecret)
		if err != nil {
			logUnauthorized(ctx, "invalid or expired token")
			apiresponse.UnAuthorized(ctx, nil, "invalid or expired token")
			ctx.Abort()
			return
		}

		ctx.Set(UserIDContextKey, claims.UserID)
		ctx.Next()
	}
}

func logUnauthorized(ctx *gin.Context, reason string) {
	logger.Warn(
		"unauthorized request",
		zap.String("reason", reason),
		zap.String("method", ctx.Request.Method),
		zap.String("path", ctx.Request.URL.Path),
		zap.String("client_ip", ctx.ClientIP()),
	)
}
