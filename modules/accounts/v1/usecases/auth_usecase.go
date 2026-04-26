package usecases

import (
	stderrors "errors"
	"strings"

	"gorm.io/gorm"

	accountErrors "api-task-management-system/modules/accounts/v1/errors"
	"api-task-management-system/modules/accounts/v1/models/users"
	"api-task-management-system/modules/accounts/v1/repositories"
	"api-task-management-system/pkg/helpers"
)

type AuthUsecase struct {
	userRepository *repositories.UserRepository
	jwtSecret      string
	jwtExpiresHour int
}

func NewAuthUsecase(userRepository *repositories.UserRepository, jwtSecret string, jwtExpiresHour int) *AuthUsecase {
	return &AuthUsecase{
		userRepository: userRepository,
		jwtSecret:      jwtSecret,
		jwtExpiresHour: jwtExpiresHour,
	}
}

func (u *AuthUsecase) Register(input users.RegisterInput) (*users.UserResponse, error) {
	email := strings.ToLower(strings.TrimSpace(input.Email))

	existingUser, err := u.userRepository.FindByEmail(email)
	if err != nil && !stderrors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	if existingUser != nil {
		return nil, accountErrors.ErrEmailAlreadyExists
	}

	passwordHash, err := helpers.HashPassword(input.Password)
	if err != nil {
		return nil, err
	}

	user := users.User{
		UUID:         helpers.NewUUID(),
		Name:         strings.TrimSpace(input.Name),
		Email:        email,
		PasswordHash: passwordHash,
	}

	if err := u.userRepository.Create(&user); err != nil {
		return nil, err
	}

	response := users.NewUserResponse(&user)
	return &response, nil
}

func (u *AuthUsecase) Login(input users.LoginInput) (*users.LoginResponse, error) {
	email := strings.ToLower(strings.TrimSpace(input.Email))

	user, err := u.userRepository.FindByEmail(email)
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
		return nil, err
	}

	return &users.LoginResponse{
		Token: token,
		User:  users.NewUserResponse(user),
	}, nil
}
