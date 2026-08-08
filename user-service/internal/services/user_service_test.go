package services

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"

	e "user-service/internal/errors"
	"user-service/internal/models"
	"user-service/internal/repository"
)

// ───────────────────────────────────────────
// GetByIDs
// ───────────────────────────────────────────

func TestGetByIDs_Success(t *testing.T) {
	mockRepo := new(repository.MockUserRepository)

	user := &models.User{
		ID:       1,
		IsActive: true,
		Role:     models.RoleUser,
	}
	mockRepo.On("GetByID", uint(1)).Return(user, nil)

	svc := NewUserService(mockRepo)
	result, err := svc.GetByIDs(1)

	assert.NoError(t, err)
	assert.Equal(t, user, result)
	mockRepo.AssertExpectations(t)
}

func TestGetByIDs_UserNotFound(t *testing.T) {
	mockRepo := new(repository.MockUserRepository)
	mockRepo.On("GetByID", uint(99)).Return(nil, errors.New("not found"))

	svc := NewUserService(mockRepo)
	result, err := svc.GetByIDs(99)

	assert.Nil(t, result)
	assert.Equal(t, e.ErrUserNotFound, err)
	mockRepo.AssertExpectations(t)
}

func TestGetByIDs_UserInactive(t *testing.T) {
	mockRepo := new(repository.MockUserRepository)

	user := &models.User{
		ID:       2,
		IsActive: false,
	}
	mockRepo.On("GetByID", uint(2)).Return(user, nil)

	svc := NewUserService(mockRepo)
	result, err := svc.GetByIDs(2)

	assert.Nil(t, result)
	assert.Equal(t, e.ErrUserInactive, err)
	mockRepo.AssertExpectations(t)
}

// ───────────────────────────────────────────
// UpdateProfile
// ───────────────────────────────────────────

func TestUpdateProfile_Success(t *testing.T) {
	mockRepo := new(repository.MockUserRepository)

	user := &models.User{ID: 1, FirstName: "Old", LastName: "Name", IsActive: true}
	mockRepo.On("GetByID", uint(1)).Return(user, nil)
	mockRepo.On("Update", user).Return(nil)

	svc := NewUserService(mockRepo)
	result, err := svc.UpdateProfile(1, "Vakha", "Selmurзаев")

	assert.NoError(t, err)
	assert.Equal(t, "Vakha", result.FirstName)
	assert.Equal(t, "Selmurзаев", result.LastName)
	mockRepo.AssertExpectations(t)
}

func TestUpdateProfile_UserNotFound(t *testing.T) {
	mockRepo := new(repository.MockUserRepository)
	mockRepo.On("GetByID", uint(99)).Return(nil, errors.New("not found"))

	svc := NewUserService(mockRepo)
	result, err := svc.UpdateProfile(99, "Name", "Last")

	assert.Nil(t, result)
	assert.Equal(t, e.ErrUserNotFound, err)
	mockRepo.AssertExpectations(t)
}

// ───────────────────────────────────────────
// BecomeOrganizer
// ───────────────────────────────────────────

func TestBecomeOrganizer_Success(t *testing.T) {
	mockRepo := new(repository.MockUserRepository)

	user := &models.User{ID: 1, Role: models.RoleUser}
	mockRepo.On("GetByID", uint(1)).Return(user, nil)
	mockRepo.On("Update", user).Return(nil)

	svc := NewUserService(mockRepo)
	result, err := svc.BecomeOrganizer(1)

	assert.NoError(t, err)
	assert.Equal(t, models.RoleOrganizer, result.Role)
	mockRepo.AssertExpectations(t)
}

func TestBecomeOrganizer_AlreadyOrganizer(t *testing.T) {
	mockRepo := new(repository.MockUserRepository)

	user := &models.User{ID: 1, Role: models.RoleOrganizer}
	mockRepo.On("GetByID", uint(1)).Return(user, nil)

	svc := NewUserService(mockRepo)
	result, err := svc.BecomeOrganizer(1)

	assert.Nil(t, result)
	assert.Equal(t, e.ErrAlreadyOrganizer, err)
	mockRepo.AssertExpectations(t)
}

func TestBecomeOrganizer_UserNotFound(t *testing.T) {
	mockRepo := new(repository.MockUserRepository)
	mockRepo.On("GetByID", uint(99)).Return(nil, errors.New("not found"))

	svc := NewUserService(mockRepo)
	result, err := svc.BecomeOrganizer(99)

	assert.Nil(t, result)
	assert.Equal(t, e.ErrUserNotFound, err)
	mockRepo.AssertExpectations(t)
}