package repositories

import (
	"context"
	"errors"

	"go.uber.org/zap"
	"gorm.io/gorm"

	"api-task-management-system/modules/accounts/v1/models/users"
)

type UserRepository struct {
	db     *gorm.DB
	logger *zap.Logger
}

func NewUserRepository(db *gorm.DB, logger *zap.Logger) *UserRepository {
	if logger == nil {
		logger = zap.NewNop()
	}

	return &UserRepository{
		db:     db,
		logger: logger,
	}
}

func (r *UserRepository) getDB(tx *gorm.DB) *gorm.DB {
	if tx != nil {
		return tx
	}

	return r.db
}

func (r *UserRepository) Create(ctx context.Context, tx *gorm.DB, user *users.User) error {
	db := r.getDB(tx).WithContext(ctx)

	if err := db.Create(user).Error; err != nil {
		r.logger.Error(
			"failed to create user",
			zap.String("email", user.Email),
			zap.String("uuid", user.UUID.String()),
			zap.Error(err),
		)
		return err
	}

	return nil
}

func (r *UserRepository) FindByEmail(ctx context.Context, tx *gorm.DB, email string) (*users.User, error) {
	db := r.getDB(tx).WithContext(ctx)

	var user users.User
	err := db.Where("email = ?", email).First(&user).Error
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			r.logger.Error(
				"failed to find user by email",
				zap.String("email", email),
				zap.Error(err),
			)
		}

		return nil, err
	}

	return &user, nil
}
