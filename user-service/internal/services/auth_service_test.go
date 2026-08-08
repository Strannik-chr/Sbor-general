package services

import (
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"

	e "user-service/internal/errors"
	"user-service/internal/models"
	"user-service/internal/repository"
	"user-service/internal/utils"
)

func newTestTokenManager() *utils.TokenManager {
	return utils.NewTokenManager("test-secret-key", time.Hour, time.Hour*24, "test-issuer")
}

// Register

func TestRegister_Success(t *testing.T) {
	mockRepo := new(repository.MockUserRepository)
	tokenManager := newTestTokenManager()

	mockRepo.On("GetByEmail", "test@mail.ru").Return(nil, errors.New("not found"))
	mockRepo.On("Create", models.AnyUser()).Return(nil)

	svc := NewAuthService(mockRepo, tokenManager)
	user, accessToken, refreshToken, err := svc.Register(
		"test@mail.ru", "password123", "Vakha", "Selmurзаев",
	)

	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.NotEmpty(t, accessToken)
	assert.NotEmpty(t, refreshToken)
	assert.Equal(t, "test@mail.ru", user.Email)
	assert.Equal(t, models.RoleUser, user.Role)
	mockRepo.AssertExpectations(t)
}

func TestRegister_EmailAlreadyExists(t *testing.T) {
	mockRepo := new(repository.MockUserRepository)
	tokenManager := newTestTokenManager()

	existingUser := &models.User{Email: "test@mail.ru"}
	mockRepo.On("GetByEmail", "test@mail.ru").Return(existingUser, nil)

	svc := NewAuthService(mockRepo, tokenManager)
	user, accessToken, refreshToken, err := svc.Register(
		"test@mail.ru", "password123", "Vakha", "Selmurзаев",
	)

	assert.Nil(t, user)
	assert.Empty(t, accessToken)
	assert.Empty(t, refreshToken)
	assert.Equal(t, e.ErrEmailAlreadyExists, err)
	mockRepo.AssertExpectations(t)
}

// Login

func TestLogin_Success(t *testing.T) {
	mockRepo := new(repository.MockUserRepository)
	tokenManager := newTestTokenManager()

	hash, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	user := &models.User{
		ID:           1,
		Email:        "test@mail.ru",
		PasswordHash: string(hash),
		IsActive:     true,
		Role:         models.RoleUser,
	}
	mockRepo.On("GetByEmail", "test@mail.ru").Return(user, nil)

	svc := NewAuthService(mockRepo, tokenManager)
	result, accessToken, refreshToken, err := svc.Login("test@mail.ru", "password123")

	assert.NoError(t, err)
	assert.Equal(t, user, result)
	assert.NotEmpty(t, accessToken)
	assert.NotEmpty(t, refreshToken)
	mockRepo.AssertExpectations(t)
}

func TestLogin_WrongPassword(t *testing.T) {
	mockRepo := new(repository.MockUserRepository)
	tokenManager := newTestTokenManager()

	hash, _ := bcrypt.GenerateFromPassword([]byte("correctpassword"), bcrypt.DefaultCost)
	user := &models.User{
		ID:           1,
		Email:        "test@mail.ru",
		PasswordHash: string(hash),
		IsActive:     true,
	}
	mockRepo.On("GetByEmail", "test@mail.ru").Return(user, nil)

	svc := NewAuthService(mockRepo, tokenManager)
	result, accessToken, refreshToken, err := svc.Login("test@mail.ru", "wrongpassword")

	assert.Nil(t, result)
	assert.Empty(t, accessToken)
	assert.Empty(t, refreshToken)
	assert.Equal(t, e.ErrInvalidCredentials, err)
	mockRepo.AssertExpectations(t)
}

func TestLogin_UserNotFound(t *testing.T) {
	mockRepo := new(repository.MockUserRepository)
	tokenManager := newTestTokenManager()

	mockRepo.On("GetByEmail", "nobody@mail.ru").Return(nil, errors.New("not found"))

	svc := NewAuthService(mockRepo, tokenManager)
	result, accessToken, refreshToken, err := svc.Login("nobody@mail.ru", "password123")

	assert.Nil(t, result)
	assert.Empty(t, accessToken)
	assert.Empty(t, refreshToken)
	assert.Equal(t, e.ErrInvalidCredentials, err)
	mockRepo.AssertExpectations(t)
}

func TestLogin_UserInactive(t *testing.T) {
	mockRepo := new(repository.MockUserRepository)
	tokenManager := newTestTokenManager()

	hash, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	user := &models.User{
		ID:           1,
		Email:        "test@mail.ru",
		PasswordHash: string(hash),
		IsActive:     false,
	}
	mockRepo.On("GetByEmail", "test@mail.ru").Return(user, nil)

	svc := NewAuthService(mockRepo, tokenManager)
	result, _, _, err := svc.Login("test@mail.ru", "password123")

	assert.Nil(t, result)
	assert.Equal(t, e.ErrUserInactive, err)
	mockRepo.AssertExpectations(t)
}