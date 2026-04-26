package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"

	"api-task-management-system/pkg/apiresponse"
	"api-task-management-system/pkg/helpers"
)

const UserIDContextKey = "user_id"

func Auth(jwtSecret string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authHeader := ctx.GetHeader("Authorization")
		if authHeader == "" {
			apiresponse.UnAuthorized(ctx, nil, "authorization header is required")
			ctx.Abort()
			return
		}

		parts := strings.Fields(authHeader)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
			apiresponse.UnAuthorized(ctx, nil, "authorization header must use Bearer token")
			ctx.Abort()
			return
		}

		claims, err := helpers.ParseJWT(parts[1], jwtSecret)
		if err != nil {
			apiresponse.UnAuthorized(ctx, nil, "invalid or expired token")
			ctx.Abort()
			return
		}

		ctx.Set(UserIDContextKey, claims.UserID)
		ctx.Next()
	}
}
