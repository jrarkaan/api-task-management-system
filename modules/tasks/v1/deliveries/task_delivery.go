package deliveries

import (
	stderrors "errors"

	"github.com/gin-gonic/gin"

	"api-task-management-system/app/middleware"
	taskErrors "api-task-management-system/modules/tasks/v1/errors"
	"api-task-management-system/modules/tasks/v1/models/tasks"
	"api-task-management-system/modules/tasks/v1/usecases"
	"api-task-management-system/pkg/apiresponse"
	"api-task-management-system/pkg/xvalidator"
)

type TaskDelivery struct {
	taskUsecase *usecases.TaskUsecase
}

func NewTaskDelivery(taskUsecase *usecases.TaskUsecase) *TaskDelivery {
	return &TaskDelivery{taskUsecase: taskUsecase}
}

func (d *TaskDelivery) List(ctx *gin.Context) {
	var query tasks.ListTaskQuery
	if err := ctx.ShouldBindQuery(&query); err != nil {
		apiresponse.BadRequest(ctx, "invalid query parameter")
		return
	}

	if err := xvalidator.Validate(query); err != nil {
		apiresponse.BadRequest(ctx, err.Error())
		return
	}

	userID, ok := getUserID(ctx)
	if !ok {
		apiresponse.UnAuthorized(ctx, nil, "unauthorized")
		return
	}

	response, err := d.taskUsecase.List(userID, query.Status)
	if err != nil {
		if stderrors.Is(err, taskErrors.ErrInvalidStatus) {
			apiresponse.BadRequest(ctx, err.Error())
			return
		}

		apiresponse.ServerError(ctx, err.Error())
		return
	}

	apiresponse.Success(ctx, nil, response)
}

func (d *TaskDelivery) Create(ctx *gin.Context) {
	var input tasks.CreateTaskInput
	if err := ctx.ShouldBindJSON(&input); err != nil {
		apiresponse.BadRequest(ctx, "invalid request body")
		return
	}

	if err := xvalidator.Validate(input); err != nil {
		apiresponse.BadRequest(ctx, err.Error())
		return
	}

	userID, ok := getUserID(ctx)
	if !ok {
		apiresponse.UnAuthorized(ctx, nil, "unauthorized")
		return
	}

	response, err := d.taskUsecase.Create(userID, input)
	if err != nil {
		apiresponse.ServerError(ctx, err.Error())
		return
	}

	apiresponse.Created(ctx, nil, response, "Created successfully")
}

func (d *TaskDelivery) Update(ctx *gin.Context) {
	var input tasks.UpdateTaskInput
	if err := ctx.ShouldBindJSON(&input); err != nil {
		apiresponse.BadRequest(ctx, "invalid request body")
		return
	}

	if err := xvalidator.Validate(input); err != nil {
		apiresponse.BadRequest(ctx, err.Error())
		return
	}

	userID, ok := getUserID(ctx)
	if !ok {
		apiresponse.UnAuthorized(ctx, nil, "unauthorized")
		return
	}

	response, err := d.taskUsecase.Update(userID, ctx.Param("id"), input)
	if err != nil {
		if stderrors.Is(err, taskErrors.ErrTaskNotFound) {
			apiresponse.NotFound(ctx, nil, err.Error())
			return
		}

		apiresponse.ServerError(ctx, err.Error())
		return
	}

	apiresponse.Success(ctx, nil, response)
}

func (d *TaskDelivery) Delete(ctx *gin.Context) {
	userID, ok := getUserID(ctx)
	if !ok {
		apiresponse.UnAuthorized(ctx, nil, "unauthorized")
		return
	}

	if err := d.taskUsecase.Delete(userID, ctx.Param("id")); err != nil {
		if stderrors.Is(err, taskErrors.ErrTaskNotFound) {
			apiresponse.NotFound(ctx, nil, err.Error())
			return
		}

		apiresponse.ServerError(ctx, err.Error())
		return
	}

	apiresponse.StatusOK(ctx, nil, "Deleted successfully")
}

func getUserID(ctx *gin.Context) (uint64, bool) {
	value, exists := ctx.Get(middleware.UserIDContextKey)
	if !exists {
		return 0, false
	}

	userID, ok := value.(uint64)
	return userID, ok
}
