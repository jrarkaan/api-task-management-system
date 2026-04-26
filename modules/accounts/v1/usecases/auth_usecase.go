package usecases

import (
	"context"
	stderrors "errors"
	"strings"

	"go.uber.org/zap"
	"gorm.io/gorm"

	accountErrors "api-task-management-system/modules/accounts/v1/errors"
	"api-task-management-system/modules/accounts/v1/models/users"
	"api-task-management-system/modules/accounts/v1/repositories"
	dbpkg "api-task-management-system/pkg/db"
	"api-task-management-system/pkg/helpers"
)

type AuthUsecase struct {
	userRepository *repositories.UserRepository
	txManager      *dbpkg.TransactionManager
	logger         *zap.Logger
	jwtSecret      string
	jwtExpiresHour int
}

func NewAuthUsecase(userRepository *repositories.UserRepository, txManager *dbpkg.TransactionManager, logger *zap.Logger, jwtSecret string, jwtExpiresHour int) *AuthUsecase {
	if logger == nil {
		logger = zap.NewNop()
	}

	return &AuthUsecase{
		userRepository: userRepository,
		txManager:      txManager,
		logger:         logger,
		jwtSecret:      jwtSecret,
		jwtExpiresHour: jwtExpiresHour,
	}
}

func (u *AuthUsecase) Register(ctx context.Context, input users.RegisterInput) (*users.UserResponse, error) {
	email := strings.ToLower(strings.TrimSpace(input.Email))

	passwordHash, err := helpers.HashPassword(input.Password)
	if err != nil {
		u.logger.Error("failed to hash password", zap.Error(err))
		return nil, err
	}

	user := users.User{
		UUID:         helpers.NewUUID(),
		Name:         strings.TrimSpace(input.Name),
		Email:        email,
		PasswordHash: passwordHash,
	}

	run := func(tx *gorm.DB) error {
		existingUser, err := u.userRepository.FindByEmail(ctx, tx, email)
		if err != nil && !stderrors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}
		if existingUser != nil {
			return accountErrors.ErrEmailAlreadyExists
		}

		return u.userRepository.Create(ctx, tx, &user)
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

	response := users.NewUserResponse(&user)
	return &response, nil
}

func (u *AuthUsecase) Login(ctx context.Context, input users.LoginInput) (*users.LoginResponse, error) {
	email := strings.ToLower(strings.TrimSpace(input.Email))

	user, err := u.userRepository.FindByEmail(ctx, nil, email)
	if err != nil {
		if stderrors.Is(err, gorm.ErrRecordNotFound) {
			return nil, accountErrors.ErrInvalidCredentials
		}

		return nil, err
	}

	if !helpers.CheckPassword(input.Password, user.PasswordHash) {
		return nil, accountErrors.ErrInvalidCredentials
	}

	token, err := helpers.GenerateJWT(user.ID, u.jwtSecret, u.jwtExpiresHour)
	if err != nil {
		u.logger.Error("failed to generate jwt", zap.Error(err), zap.Uint64("user_id", user.ID))
		return nil, err
	}

	return &users.LoginResponse{
		Token: token,
		User:  users.NewUserResponse(user),
	}, nil
}
