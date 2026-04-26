package deliveries

import (
	stderrors "errors"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"api-task-management-system/app/middleware"
	taskErrors "api-task-management-system/modules/tasks/v1/errors"
	"api-task-management-system/modules/tasks/v1/models/tasks"
	"api-task-management-system/modules/tasks/v1/usecases"
	"api-task-management-system/pkg/apiresponse"
	"api-task-management-system/pkg/logger"
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
		logger.Warn("invalid task list query", zap.Error(err))
		apiresponse.BadRequest(ctx, "invalid query parameter")
		return
	}

	if err := xvalidator.Validate(query); err != nil {
		logger.Warn("task list validation failed", zap.Error(err))
		apiresponse.BadRequest(ctx, err.Error())
		return
	}

	userID, ok := getUserID(ctx)
	if !ok {
		logger.Warn("task list missing user context")
		apiresponse.UnAuthorized(ctx, nil, "unauthorized")
		return
	}

	response, err := d.taskUsecase.List(userID, query.Status)
	if err != nil {
		if stderrors.Is(err, taskErrors.ErrInvalidStatus) {
			logger.Warn("task list invalid status", zap.Error(err), zap.Uint64("user_id", userID))
			apiresponse.BadRequest(ctx, err.Error())
			return
		}

		logger.Error("task list failed", zap.Error(err), zap.Uint64("user_id", userID))
		apiresponse.ServerError(ctx, err.Error())
		return
	}

	apiresponse.Success(ctx, nil, response)
}

func (d *TaskDelivery) Create(ctx *gin.Context) {
	var input tasks.CreateTaskInput
	if err := ctx.ShouldBindJSON(&input); err != nil {
		logger.Warn("invalid create task request body", zap.Error(err))
		apiresponse.BadRequest(ctx, "invalid request body")
		return
	}

	if err := xvalidator.Validate(input); err != nil {
		logger.Warn("create task validation failed", zap.Error(err))
		apiresponse.BadRequest(ctx, err.Error())
		return
	}

	userID, ok := getUserID(ctx)
	if !ok {
		logger.Warn("create task missing user context")
		apiresponse.UnAuthorized(ctx, nil, "unauthorized")
		return
	}

	response, err := d.taskUsecase.Create(userID, input)
	if err != nil {
		logger.Error("create task failed", zap.Error(err), zap.Uint64("user_id", userID))
		apiresponse.ServerError(ctx, err.Error())
		return
	}

	apiresponse.Created(ctx, nil, response, "Created successfully")
}

func (d *TaskDelivery) Update(ctx *gin.Context) {
	var input tasks.UpdateTaskInput
	if err := ctx.ShouldBindJSON(&input); err != nil {
		logger.Warn("invalid update task request body", zap.Error(err))
		apiresponse.BadRequest(ctx, "invalid request body")
		return
	}

	if err := xvalidator.Validate(input); err != nil {
		logger.Warn("update task validation failed", zap.Error(err))
		apiresponse.BadRequest(ctx, err.Error())
		return
	}

	userID, ok := getUserID(ctx)
	if !ok {
		logger.Warn("update task missing user context")
		apiresponse.UnAuthorized(ctx, nil, "unauthorized")
		return
	}

	response, err := d.taskUsecase.Update(userID, ctx.Param("id"), input)
	if err != nil {
		if stderrors.Is(err, taskErrors.ErrTaskNotFound) {
			logger.Warn("task update not found", zap.String("task_id", ctx.Param("id")), zap.Uint64("user_id", userID))
			apiresponse.NotFound(ctx, nil, err.Error())
			return
		}

		logger.Error("task update failed", zap.Error(err), zap.String("task_id", ctx.Param("id")), zap.Uint64("user_id", userID))
		apiresponse.ServerError(ctx, err.Error())
		return
	}

	apiresponse.Success(ctx, nil, response)
}

func (d *TaskDelivery) Delete(ctx *gin.Context) {
	userID, ok := getUserID(ctx)
	if !ok {
		logger.Warn("delete task missing user context")
		apiresponse.UnAuthorized(ctx, nil, "unauthorized")
		return
	}

	if err := d.taskUsecase.Delete(userID, ctx.Param("id")); err != nil {
		if stderrors.Is(err, taskErrors.ErrTaskNotFound) {
			logger.Warn("task delete not found", zap.String("task_id", ctx.Param("id")), zap.Uint64("user_id", userID))
			apiresponse.NotFound(ctx, nil, err.Error())
			return
		}

		logger.Error("task delete failed", zap.Error(err), zap.String("task_id", ctx.Param("id")), zap.Uint64("user_id", userID))
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
