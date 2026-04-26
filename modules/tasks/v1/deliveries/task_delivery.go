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

// List godoc
// @Summary      List tasks
// @Description  Get a list of user tasks with optional pagination and status filter.
// @Tags         Tasks
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        status query string false "Filter by status" Enums(pending, in-progress, done)
// @Param        page query int false "Page number" default(1)
// @Param        limit query int false "Items per page (max 100)" default(10)
// @Success      200  {object}  apiresponse.SwaggerTaskListResponse
// @Failure      400  {object}  apiresponse.SwaggerErrorResponse
// @Failure      401  {object}  apiresponse.SwaggerErrorResponse
// @Failure      500  {object}  apiresponse.SwaggerErrorResponse
// @Router       /tasks [get]
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

	response, meta, err := d.taskUsecase.List(ctx.Request.Context(), userID, query.Status, query.Page, query.Limit)
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

	apiresponse.Success(ctx, meta, response)
}

// Create godoc
// @Summary      Create a new task
// @Description  Creates a task under the authenticated user.
// @Tags         Tasks
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request body tasks.CreateTaskInput true "Task creation details"
// @Success      201  {object}  apiresponse.SwaggerTaskResponse
// @Failure      400  {object}  apiresponse.SwaggerErrorResponse
// @Failure      401  {object}  apiresponse.SwaggerErrorResponse
// @Failure      500  {object}  apiresponse.SwaggerErrorResponse
// @Router       /tasks [post]
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

	response, err := d.taskUsecase.Create(ctx.Request.Context(), userID, input)
	if err != nil {
		logger.Error("create task failed", zap.Error(err), zap.Uint64("user_id", userID))
		apiresponse.ServerError(ctx, err.Error())
		return
	}

	apiresponse.Created(ctx, nil, response, "Created successfully")
}

// Update godoc
// @Summary      Update an existing task
// @Description  Partially updates a task's fields. At least one field must be provided. Omitted fields are left unchanged.
// @Tags         Tasks
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id path string true "Task UUID"
// @Param        request body tasks.UpdateTaskInput true "Partial task update fields"
// @Success      200  {object}  apiresponse.SwaggerTaskResponse
// @Failure      400  {object}  apiresponse.SwaggerErrorResponse
// @Failure      401  {object}  apiresponse.SwaggerErrorResponse
// @Failure      404  {object}  apiresponse.SwaggerErrorResponse
// @Failure      500  {object}  apiresponse.SwaggerErrorResponse
// @Router       /tasks/{id} [put]
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

	if !input.HasUpdateFields() {
		logger.Warn("update task missing fields")
		apiresponse.BadRequest(ctx, "at least one field must be provided")
		return
	}

	userID, ok := getUserID(ctx)
	if !ok {
		logger.Warn("update task missing user context")
		apiresponse.UnAuthorized(ctx, nil, "unauthorized")
		return
	}

	response, err := d.taskUsecase.Update(ctx.Request.Context(), userID, ctx.Param("id"), input)
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

// Delete godoc
// @Summary      Delete a task
// @Description  Deletes a given task belonging to the authenticated user.
// @Tags         Tasks
// @Produce      json
// @Security     BearerAuth
// @Param        id path string true "Task UUID"
// @Success      200  {object}  apiresponse.SwaggerSuccessResponse
// @Failure      401  {object}  apiresponse.SwaggerErrorResponse
// @Failure      404  {object}  apiresponse.SwaggerErrorResponse
// @Failure      500  {object}  apiresponse.SwaggerErrorResponse
// @Router       /tasks/{id} [delete]
func (d *TaskDelivery) Delete(ctx *gin.Context) {
	userID, ok := getUserID(ctx)
	if !ok {
		logger.Warn("delete task missing user context")
		apiresponse.UnAuthorized(ctx, nil, "unauthorized")
		return
	}

	if err := d.taskUsecase.Delete(ctx.Request.Context(), userID, ctx.Param("id")); err != nil {
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
