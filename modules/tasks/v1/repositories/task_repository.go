package repositories

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"go.uber.org/zap"
	"gorm.io/gorm"

	"api-task-management-system/modules/tasks/v1/models/tasks"
	"api-task-management-system/pkg/pagination"
)

type TaskRepository struct {
	db     *gorm.DB
	logger *zap.Logger
}

func NewTaskRepository(db *gorm.DB, logger *zap.Logger) *TaskRepository {
	if logger == nil {
		logger = zap.NewNop()
	}

	return &TaskRepository{
		db:     db,
		logger: logger,
	}
}

func (r *TaskRepository) getDB(tx *gorm.DB) *gorm.DB {
	if tx != nil {
		return tx
	}

	return r.db
}

func (r *TaskRepository) Create(ctx context.Context, tx *gorm.DB, task *tasks.Task) error {
	db := r.getDB(tx).WithContext(ctx)

	if err := db.Create(task).Error; err != nil {
		r.logger.Error(
			"failed to create task",
			zap.String("uuid", task.UUID.String()),
			zap.Uint64("user_id", task.UserID),
			zap.String("status", task.Status),
			zap.Error(err),
		)
		return err
	}

	return nil
}

func (r *TaskRepository) ListByUser(ctx context.Context, tx *gorm.DB, userID uint64, status string, page int, limit int) ([]tasks.Task, int64, error) {
	db := r.getDB(tx).WithContext(ctx)

	var taskList []tasks.Task
	baseQuery := db.Model(&tasks.Task{}).Where("user_id = ?", userID)
	if status != "" {
		baseQuery = baseQuery.Where("status = ?", status)
	}

	var totalRows int64
	if err := baseQuery.Count(&totalRows).Error; err != nil {
		r.logger.Error(
			"failed to count tasks",
			zap.Uint64("user_id", userID),
			zap.String("status", status),
			zap.Error(err),
		)
		return nil, 0, err
	}

	err := baseQuery.
		Order("created_at DESC").
		Limit(limit).
		Offset(pagination.Offset(page, limit)).
		Find(&taskList).Error

	if err != nil {
		r.logger.Error(
			"failed to list tasks",
			zap.Uint64("user_id", userID),
			zap.String("status", status),
			zap.Error(err),
		)
		return nil, 0, err
	}

	return taskList, totalRows, nil
}

func (r *TaskRepository) FindByUUIDAndUser(ctx context.Context, tx *gorm.DB, taskUUID uuid.UUID, userID uint64) (*tasks.Task, error) {
	db := r.getDB(tx).WithContext(ctx)

	var task tasks.Task
	err := db.Where("uuid = ? AND user_id = ?", taskUUID, userID).First(&task).Error
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			r.logger.Error(
				"failed to find task by uuid and user",
				zap.String("uuid", taskUUID.String()),
				zap.Uint64("user_id", userID),
				zap.Error(err),
			)
		}

		return nil, err
	}

	return &task, nil
}

func (r *TaskRepository) Update(ctx context.Context, tx *gorm.DB, task *tasks.Task) error {
	db := r.getDB(tx).WithContext(ctx)

	if err := db.Save(task).Error; err != nil {
		r.logger.Error(
			"failed to update task",
			zap.Uint64("task_id", task.ID),
			zap.String("uuid", task.UUID.String()),
			zap.Uint64("user_id", task.UserID),
			zap.String("status", task.Status),
			zap.Error(err),
		)
		return err
	}

	return nil
}

func (r *TaskRepository) Delete(ctx context.Context, tx *gorm.DB, task *tasks.Task) error {
	db := r.getDB(tx).WithContext(ctx)

	if err := db.Delete(task).Error; err != nil {
		r.logger.Error(
			"failed to delete task",
			zap.Uint64("task_id", task.ID),
			zap.String("uuid", task.UUID.String()),
			zap.Uint64("user_id", task.UserID),
			zap.String("status", task.Status),
			zap.Error(err),
		)
		return err
	}

	return nil
}
