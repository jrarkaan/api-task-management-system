package usecases

import (
	"context"
	stderrors "errors"
	"strings"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"
	"gorm.io/gorm"

	taskErrors "api-task-management-system/modules/tasks/v1/errors"
	"api-task-management-system/modules/tasks/v1/models/tasks"
	"api-task-management-system/modules/tasks/v1/repositories"
	dbpkg "api-task-management-system/pkg/db"
	"api-task-management-system/pkg/helpers"
)

type TaskUsecase struct {
	taskRepository *repositories.TaskRepository
	txManager      *dbpkg.TransactionManager
	logger         *zap.Logger
}

func NewTaskUsecase(taskRepository *repositories.TaskRepository, txManager *dbpkg.TransactionManager, logger *zap.Logger) *TaskUsecase {
	if logger == nil {
		logger = zap.NewNop()
	}

	return &TaskUsecase{
		taskRepository: taskRepository,
		txManager:      txManager,
		logger:         logger,
	}
}

func (u *TaskUsecase) List(ctx context.Context, userID uint64, status string) ([]tasks.TaskResponse, error) {
	if status != "" && !tasks.IsValidStatus(status) {
		return nil, taskErrors.ErrInvalidStatus
	}

	taskList, err := u.taskRepository.ListByUser(ctx, nil, userID, status)
	if err != nil {
		return nil, err
	}

	return tasks.NewTaskResponses(taskList), nil
}

func (u *TaskUsecase) Create(ctx context.Context, userID uint64, input tasks.CreateTaskInput) (*tasks.TaskResponse, error) {
	deadline, err := parseDeadline(input.Deadline)
	if err != nil {
		return nil, err
	}

	status := input.Status
	if status == "" {
		status = tasks.StatusPending
	}

	description := nullableString(input.Description)
	task := tasks.Task{
		UUID:        helpers.NewUUID(),
		UserID:      userID,
		Title:       strings.TrimSpace(input.Title),
		Description: description,
		Status:      status,
		Deadline:    deadline,
	}

	if err := u.taskRepository.Create(ctx, nil, &task); err != nil {
		return nil, err
	}

	response := tasks.NewTaskResponse(&task)
	return &response, nil
}

func (u *TaskUsecase) Update(ctx context.Context, userID uint64, taskID string, input tasks.UpdateTaskInput) (*tasks.TaskResponse, error) {
	taskUUID, err := uuid.Parse(taskID)
	if err != nil {
		return nil, taskErrors.ErrTaskNotFound
	}

	deadline, err := parseDeadline(input.Deadline)
	if err != nil {
		return nil, err
	}

	var task *tasks.Task
	run := func(tx *gorm.DB) error {
		var err error
		task, err = u.taskRepository.FindByUUIDAndUser(ctx, tx, taskUUID, userID)
		if err != nil {
			if stderrors.Is(err, gorm.ErrRecordNotFound) {
				return taskErrors.ErrTaskNotFound
			}

			return err
		}

		task.Title = strings.TrimSpace(input.Title)
		task.Description = nullableString(input.Description)
		if input.Status != "" {
			task.Status = input.Status
		}
		if input.Deadline != "" {
			task.Deadline = deadline
		}

		return u.taskRepository.Update(ctx, tx, task)
	}

	if u.txManager != nil {
		if err := u.txManager.WithTransaction(ctx, run); err != nil {
			return nil, err
		}
	} else {
		if err := run(nil); err != nil {
			return nil, err
		}
	}

	response := tasks.NewTaskResponse(task)
	return &response, nil
}

func (u *TaskUsecase) Delete(ctx context.Context, userID uint64, taskID string) error {
	taskUUID, err := uuid.Parse(taskID)
	if err != nil {
		return taskErrors.ErrTaskNotFound
	}

	run := func(tx *gorm.DB) error {
		task, err := u.taskRepository.FindByUUIDAndUser(ctx, tx, taskUUID, userID)
		if err != nil {
			if stderrors.Is(err, gorm.ErrRecordNotFound) {
				return taskErrors.ErrTaskNotFound
			}

			return err
		}

		return u.taskRepository.Delete(ctx, tx, task)
	}

	if u.txManager != nil {
		return u.txManager.WithTransaction(ctx, run)
	}

	return run(nil)
}

func parseDeadline(value string) (*time.Time, error) {
	if value == "" {
		return nil, nil
	}

	deadline, err := time.Parse("2006-01-02", value)
	if err != nil {
		return nil, err
	}

	return &deadline, nil
}

func nullableString(value string) *string {
	value = strings.TrimSpace(value)
	if value == "" {
		return nil
	}

	return &value
}
