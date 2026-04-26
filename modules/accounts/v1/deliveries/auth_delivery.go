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

// Register godoc
// @Summary      Register a new user
// @Description  Creates a new user account and returns user details.
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        request body users.RegisterInput true "Register details"
// @Success      201  {object}  apiresponse.SwaggerSuccessResponse
// @Failure      400  {object}  apiresponse.SwaggerErrorResponse
// @Failure      409  {object}  apiresponse.SwaggerErrorResponse
// @Failure      500  {object}  apiresponse.SwaggerErrorResponse
// @Router       /accounts/v1/auth/register [post]
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

	response, err := d.authUsecase.Register(ctx.Request.Context(), input)
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

// Login godoc
// @Summary      User login
// @Description  Authenticates a user and returns a JWT token.
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        request body users.LoginInput true "Login credentials"
// @Success      200  {object}  apiresponse.SwaggerAuthLoginResponse
// @Failure      400  {object}  apiresponse.SwaggerErrorResponse
// @Failure      401  {object}  apiresponse.SwaggerErrorResponse
// @Failure      500  {object}  apiresponse.SwaggerErrorResponse
// @Router       /accounts/v1/auth/login [post]
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

	response, err := d.authUsecase.Login(ctx.Request.Context(), input)
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
