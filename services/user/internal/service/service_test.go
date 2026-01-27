package service

import (
	"context"
	"testing"

	"github.com/katatrina/airbnb-clone/services/user/internal/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestCreateUser_Success(t *testing.T) {
	mockUserRepo := new(MockUserRepository)
	mockTokenMaker := new(MockTokenMaker)
	userSvc := NewUserService(mockUserRepo, mockTokenMaker)

	ctx := context.Background()
	arg := model.CreateUserParams{
		Email:       "newuser@example.com",
		Password:    "securepassword123",
		DisplayName: "New User",
	}

	mockUserRepo.On("CheckEmailExists", ctx, arg.Email).Return(false, nil)
	mockUserRepo.On("CreateUser", ctx, mock.AnythingOfType("*model.User")).Return(nil)

	user, err := userSvc.CreateUser(ctx, arg)

	require.NoError(t, err)
	require.NotNil(t, user)
	assert.Equal(t, arg.Email, user.Email)
	assert.Equal(t, arg.DisplayName, user.DisplayName)
	assert.NotEmpty(t, user.ID)
	assert.NotEmpty(t, user.PasswordHash)
	assert.False(t, user.EmailVerified)

	mockUserRepo.AssertExpectations(t)
}
