package repositories

import (
	"github.com/google/uuid"
	"gorm.io/gorm"

	"api-task-management-system/modules/tasks/v1/models/tasks"
)

type TaskRepository struct {
	db *gorm.DB
}

func NewTaskRepository(db *gorm.DB) *TaskRepository {
	return &TaskRepository{db: db}
}

func (r *TaskRepository) Create(task *tasks.Task) error {
	return r.db.Create(task).Error
}

func (r *TaskRepository) ListByUser(userID uint64, status string) ([]tasks.Task, error) {
	var taskList []tasks.Task
	query := r.db.Where("user_id = ?", userID)
	if status != "" {
		query = query.Where("status = ?", status)
	}

	err := query.Order("created_at DESC").Find(&taskList).Error
	return taskList, err
}

func (r *TaskRepository) FindByUUIDAndUser(taskUUID uuid.UUID, userID uint64) (*tasks.Task, error) {
	var task tasks.Task
	err := r.db.Where("uuid = ? AND user_id = ?", taskUUID, userID).First(&task).Error
	if err != nil {
		return nil, err
	}

	return &task, nil
}

func (r *TaskRepository) Update(task *tasks.Task) error {
	return r.db.Save(task).Error
}

func (r *TaskRepository) Delete(task *tasks.Task) error {
	return r.db.Delete(task).Error
}
