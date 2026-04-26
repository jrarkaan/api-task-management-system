package deliveries

import (
	stderrors "errors"

	"github.com/gin-gonic/gin"

	accountErrors "api-task-management-system/modules/accounts/v1/errors"
	"api-task-management-system/modules/accounts/v1/models/users"
	"api-task-management-system/modules/accounts/v1/usecases"
	"api-task-management-system/pkg/apiresponse"
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
		apiresponse.BadRequest(ctx, "invalid request body")
		return
	}

	if err := xvalidator.Validate(input); err != nil {
		apiresponse.BadRequest(ctx, err.Error())
		return
	}

	response, err := d.authUsecase.Register(input)
	if err != nil {
		if stderrors.Is(err, accountErrors.ErrEmailAlreadyExists) {
			apiresponse.Conflict(ctx, nil, err.Error())
			return
		}

		apiresponse.ServerError(ctx, err.Error())
		return
	}

	apiresponse.Created(ctx, nil, response, "Created successfully")
}

func (d *AuthDelivery) Login(ctx *gin.Context) {
	var input users.LoginInput
	if err := ctx.ShouldBindJSON(&input); err != nil {
		apiresponse.BadRequest(ctx, "invalid request body")
		return
	}

	if err := xvalidator.Validate(input); err != nil {
		apiresponse.BadRequest(ctx, err.Error())
		return
	}

	response, err := d.authUsecase.Login(input)
	if err != nil {
		if stderrors.Is(err, accountErrors.ErrInvalidCredentials) {
			apiresponse.UnAuthorized(ctx, nil, err.Error())
			return
		}

		apiresponse.ServerError(ctx, err.Error())
		return
	}

	apiresponse.Success(ctx, nil, response)
}
