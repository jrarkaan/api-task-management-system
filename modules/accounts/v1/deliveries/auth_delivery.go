package deliveries

import (
	stderrors "errors"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	accountErrors "api-task-management-system/modules/accounts/v1/errors"
	"api-task-management-system/modules/accounts/v1/models/users"
	"api-task-management-system/modules/accounts/v1/usecases"
	"api-task-management-system/pkg/apiresponse"
	"api-task-management-system/pkg/logger"
	"api-task-management-system/pkg/xvalidator"
)

type AuthDelivery struct {
	authUsecase *usecases.AuthUsecase
}

func NewAuthDelivery(authUsecase *usecases.AuthUsecase) *AuthDelivery {
	return &AuthDelivery{authUsecase: authUsecase}
}

func (d *AuthDelivery) Register(ctx *gin.Context) {
	var input users.RegisterInput
	if err := ctx.ShouldBindJSON(&input); err != nil {
		logger.Warn("invalid register request body", zap.Error(err))
		apiresponse.BadRequest(ctx, "invalid request body")
		return
	}

	if err := xvalidator.Validate(input); err != nil {
		logger.Warn("register validation failed", zap.Error(err))
		apiresponse.BadRequest(ctx, err.Error())
		return
	}

	response, err := d.authUsecase.Register(input)
	if err != nil {
		if stderrors.Is(err, accountErrors.ErrEmailAlreadyExists) {
			logger.Warn("register conflict", zap.Error(err))
			apiresponse.Conflict(ctx, nil, err.Error())
			return
		}

		logger.Error("register failed", zap.Error(err))
		apiresponse.ServerError(ctx, err.Error())
		return
	}

	apiresponse.Created(ctx, nil, response, "Created successfully")
}

func (d *AuthDelivery) Login(ctx *gin.Context) {
	var input users.LoginInput
	if err := ctx.ShouldBindJSON(&input); err != nil {
		logger.Warn("invalid login request body", zap.Error(err))
		apiresponse.BadRequest(ctx, "invalid request body")
		return
	}

	if err := xvalidator.Validate(input); err != nil {
		logger.Warn("login validation failed", zap.Error(err))
		apiresponse.BadRequest(ctx, err.Error())
		return
	}

	response, err := d.authUsecase.Login(input)
	if err != nil {
		if stderrors.Is(err, accountErrors.ErrInvalidCredentials) {
			logger.Warn("login unauthorized", zap.Error(err))
			apiresponse.UnAuthorized(ctx, nil, err.Error())
			return
		}

		logger.Error("login failed", zap.Error(err))
		apiresponse.ServerError(ctx, err.Error())
		return
	}

	apiresponse.Success(ctx, nil, response)
}
