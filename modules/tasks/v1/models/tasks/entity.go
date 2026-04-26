package tasks

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

const (
	StatusPending    = "pending"
	StatusInProgress = "in-progress"
	StatusDone       = "done"
)

type Task struct {
	ID          uint64         `gorm:"primaryKey" json:"-"`
	UUID        uuid.UUID      `gorm:"type:uuid;uniqueIndex;not null" json:"uuid"`
	UserID      uint64         `gorm:"not null;index" json:"-"`
	Title       string         `gorm:"size:255;not null" json:"title"`
	Description *string        `gorm:"type:text" json:"description,omitempty"`
	Status      string         `gorm:"size:20;not null;default:pending" json:"status"`
	Deadline    *time.Time     `gorm:"type:date" json:"deadline,omitempty"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

func (Task) TableName() string {
	return "tasks"
}

func (t *Task) BeforeCreate(tx *gorm.DB) error {
	if t.UUID == uuid.Nil {
		t.UUID = uuid.New()
	}
	if t.Status == "" {
		t.Status = StatusPending
	}

	return nil
}

func IsValidStatus(status string) bool {
	return status == StatusPending || status == StatusInProgress || status == StatusDone
}
