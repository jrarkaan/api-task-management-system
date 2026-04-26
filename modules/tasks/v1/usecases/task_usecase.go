package usecases

import (
	stderrors "errors"
	"strings"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"
	"gorm.io/gorm"

	taskErrors "api-task-management-system/modules/tasks/v1/errors"
	"api-task-management-system/modules/tasks/v1/models/tasks"
	"api-task-management-system/modules/tasks/v1/repositories"
	"api-task-management-system/pkg/helpers"
	"api-task-management-system/pkg/logger"
)

type TaskUsecase struct {
	taskRepository *repositories.TaskRepository
}

func NewTaskUsecase(taskRepository *repositories.TaskRepository) *TaskUsecase {
	return &TaskUsecase{taskRepository: taskRepository}
}

func (u *TaskUsecase) List(userID uint64, status string) ([]tasks.TaskResponse, error) {
	if status != "" && !tasks.IsValidStatus(status) {
		return nil, taskErrors.ErrInvalidStatus
	}

	taskList, err := u.taskRepository.ListByUser(userID, status)
	if err != nil {
		logger.Error("failed to list tasks", zap.Error(err), zap.Uint64("user_id", userID))
		return nil, err
	}

	return tasks.NewTaskResponses(taskList), nil
}

func (u *TaskUsecase) Create(userID uint64, input tasks.CreateTaskInput) (*tasks.TaskResponse, error) {
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

	if err := u.taskRepository.Create(&task); err != nil {
		logger.Error("failed to create task", zap.Error(err), zap.Uint64("user_id", userID))
		return nil, err
	}

	response := tasks.NewTaskResponse(&task)
	return &response, nil
}

func (u *TaskUsecase) Update(userID uint64, taskID string, input tasks.UpdateTaskInput) (*tasks.TaskResponse, error) {
	taskUUID, err := uuid.Parse(taskID)
	if err != nil {
		return nil, taskErrors.ErrTaskNotFound
	}

	task, err := u.taskRepository.FindByUUIDAndUser(taskUUID, userID)
	if err != nil {
		if stderrors.Is(err, gorm.ErrRecordNotFound) {
			return nil, taskErrors.ErrTaskNotFound
		}

		logger.Error("failed to find task for update", zap.Error(err), zap.String("task_id", taskID), zap.Uint64("user_id", userID))
		return nil, err
	}

	deadline, err := parseDeadline(input.Deadline)
	if err != nil {
		return nil, err
	}

	task.Title = strings.TrimSpace(input.Title)
	task.Description = nullableString(input.Description)
	if input.Status != "" {
		task.Status = input.Status
	}
	if input.Deadline != "" {
		task.Deadline = deadline
	}

	if err := u.taskRepository.Update(task); err != nil {
		logger.Error("failed to update task", zap.Error(err), zap.String("task_id", taskID), zap.Uint64("user_id", userID))
		return nil, err
	}

	response := tasks.NewTaskResponse(task)
	return &response, nil
}

func (u *TaskUsecase) Delete(userID uint64, taskID string) error {
	taskUUID, err := uuid.Parse(taskID)
	if err != nil {
		return taskErrors.ErrTaskNotFound
	}

	task, err := u.taskRepository.FindByUUIDAndUser(taskUUID, userID)
	if err != nil {
		if stderrors.Is(err, gorm.ErrRecordNotFound) {
			return taskErrors.ErrTaskNotFound
		}

		logger.Error("failed to find task for delete", zap.Error(err), zap.String("task_id", taskID), zap.Uint64("user_id", userID))
		return err
	}

	if err := u.taskRepository.Delete(task); err != nil {
		logger.Error("failed to delete task", zap.Error(err), zap.String("task_id", taskID), zap.Uint64("user_id", userID))
		return err
	}

	return nil
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
