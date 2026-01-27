package service

import (
	"context"
	"errors"
	"testing"

	"github.com/katatrina/airbnb-clone/services/user/internal/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestCreateUser_Success(t *testing.T) {
	mockUserRepo := new(MockUserRepository)
	userSvc := NewUserService(mockUserRepo, nil)

	ctx := context.Background()
	arg := model.CreateUserParams{
		Email:       "newuser@example.com",
		Password:    "securepassword123",
		DisplayName: "New User",
	}

	mockUserRepo.On("CheckEmailExists", ctx, arg.Email).
		Return(false, nil)
	mockUserRepo.On("CreateUser", ctx, mock.AnythingOfType("*model.User")).
		Return(nil)

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

func TestCreateUser_EmailAlreadyExists(t *testing.T) {
	mockUserRepo := new(MockUserRepository)
	userSvc := NewUserService(mockUserRepo, nil)

	ctx := context.Background()
	arg := model.CreateUserParams{
		Email:       "existing@example.com",
		Password:    "password123",
		DisplayName: "User",
	}

	mockUserRepo.On("CheckEmailExists", ctx, arg.Email).
		Return(true, nil)

	user, err := userSvc.CreateUser(ctx, arg)

	require.Error(t, err)
	require.Nil(t, user)
	assert.ErrorIs(t, err, model.ErrEmailAlreadyExists)

	mockUserRepo.AssertExpectations(t)
}

func TestCreateUser_RepositoryError(t *testing.T) {
	mockUserRepo := new(MockUserRepository)
	userSvc := NewUserService(mockUserRepo, nil)

	ctx := context.Background()
	arg := model.CreateUserParams{
		Email:       "user@example.com",
		Password:    "password123",
		DisplayName: "User",
	}

	dbErr := errors.New("database connection failed")
	mockUserRepo.On("CheckEmailExists", ctx, arg.Email).
		Return(false, dbErr)

	user, err := userSvc.CreateUser(ctx, arg)

	require.Error(t, err)
	require.Nil(t, user)
	assert.ErrorIs(t, err, dbErr)

	mockUserRepo.AssertExpectations(t)

}

func TestLoginUser(t *testing.T) {
	testCases := []struct {
		name           string
		arg            model.LoginUserParams
		setupMock      func(*MockUserRepository, *MockTokenMaker)
		wantError      bool
		expectedErrMsg string
	}{
		{
			name: "successful login",
			arg: model.LoginUserParams{
				Email:    "user@example.com",
				Password: "correctpassword",
			},
			setupMock: func(userRepo *MockUserRepository, tokenMaker *MockTokenMaker) {
				userRepo.On("FindUserByEmail", mock.Anything, "user@example.com").
					Return(mock.AnythingOfType("*model.User"), nil)

				tokenMaker.On("CreateToken", "user-123").
					Return("jwt-token-xyz", nil)

				userRepo.On("UpdateUserLastLogin", mock.Anything, "user-123", mock.AnythingOfType("time.Time")).
					Return(nil)
			},
			wantError: false,
		},
		{
			name: "email not found",
			arg: model.LoginUserParams{
				Email:    "nonexistent@example.com",
				Password: "anypassword",
			},
			setupMock: func(userRepo *MockUserRepository, tokenMaker *MockTokenMaker) {
				userRepo.On("FindUserByEmail", mock.Anything, "nonexistent@example.com").
					Return(nil, model.ErrUserNotFound)
			},
			wantError:      true,
			expectedErrMsg: "incorrect email or password",
		},
		{
			name: "incorrect password",
			arg: model.LoginUserParams{
				Email:    "user@example.com",
				Password: "wrongpassword",
			},
			setupMock: func(userRepo *MockUserRepository, tokenMaker *MockTokenMaker) {
				userRepo.On("FindUserByEmail", mock.Anything, "user@example.com").
					Return(mock.AnythingOfType("*model.User"), nil)
			},
			wantError:      true,
			expectedErrMsg: "incorrect email or password",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockUserRepo := new(MockUserRepository)
			mockTokenMaker := new(MockTokenMaker)

			tc.setupMock(mockUserRepo, mockTokenMaker)

			userSvc := NewUserService(mockUserRepo, mockTokenMaker)

			result, err := userSvc.LoginUser(context.Background(), tc.arg)

			if tc.wantError {
				require.Error(t, err)
				if tc.expectedErrMsg != "" {
					assert.Contains(t, err.Error(), tc.expectedErrMsg)
				}
				require.Nil(t, result)
			} else {
				require.NoError(t, err)
				require.NotNil(t, result)
				assert.NotEmpty(t, result.AccessToken)
			}

			mockUserRepo.AssertExpectations(t)
			mockTokenMaker.AssertExpectations(t)
		})
	}
}
