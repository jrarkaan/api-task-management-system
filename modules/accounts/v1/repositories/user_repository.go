package repositories

import (
	"gorm.io/gorm"

	"api-task-management-system/modules/accounts/v1/models/users"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(user *users.User) error {
	return r.db.Create(user).Error
}

func (r *UserRepository) FindByEmail(email string) (*users.User, error) {
	var user users.User
	err := r.db.Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}

	return &user, nil
}
