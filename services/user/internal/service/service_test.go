// service_test.go
// Unit tests for UserService - focuses on business logic and mock interactions
package service

import (
	"context"
	"errors"
	"testing"

	"github.com/katatrina/airbnb-clone/services/user/internal/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)

// createTestUser creates a user with bcrypt-hashed password for testing login flow
func createTestUser(id, email, password string) *model.User {
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
	return &model.User{
		ID:           id,
		DisplayName:  "Test User",
		Email:        email,
		PasswordHash: string(hashedPassword),
	}
}

func newMocksAndService() (*MockUserRepository, *MockTokenMaker, *UserService) {
	mockRepo := new(MockUserRepository)
	mockToken := new(MockTokenMaker)
	service := NewUserService(mockRepo, mockToken)

	return mockRepo, mockToken, service
}

func TestCreateUser(t *testing.T) {
	testCases := []struct {
		name        string
		input       model.CreateUserParams
		setupMock   func(*MockUserRepository, *MockTokenMaker)
		wantErr     bool
		expectedErr error
		validate    func(t *testing.T, input model.CreateUserParams, user *model.User)
	}{
		{
			name: "success - creates user with valid input",
			input: model.CreateUserParams{
				DisplayName: "New User",
				Email:       "newuser@example.com",
				Password:    "strongPassword123",
			},
			setupMock: func(mockRepo *MockUserRepository, mockToken *MockTokenMaker) {
				mockRepo.On("CheckEmailExists", mock.Anything, "newuser@example.com").
					Return(false, nil)
				mockRepo.On("CreateUser", mock.Anything, mock.AnythingOfType("*model.User")).
					Return(nil)
			},
			wantErr:     false,
			expectedErr: nil,
			validate: func(t *testing.T, input model.CreateUserParams, user *model.User) {
				assert.Equal(t, input.Email, user.Email)
				assert.Equal(t, input.DisplayName, user.DisplayName)
				assert.NotEmpty(t, user.ID)
				assert.NotEqual(t, input.Password, user.PasswordHash)
			},
		},
		{
			name: "error - email already exists",
			input: model.CreateUserParams{
				DisplayName: "Existing User",
				Email:       "existing@example.com",
				Password:    "strongPassword123",
			},
			setupMock: func(mockRepo *MockUserRepository, mockToken *MockTokenMaker) {
				mockRepo.On("CheckEmailExists", mock.Anything, "existing@example.com").
					Return(true, nil)
			},
			wantErr:     true,
			expectedErr: model.ErrEmailAlreadyExists,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockRepo, mockToken, svc := newMocksAndService()
			tc.setupMock(mockRepo, mockToken)

			ctx := context.Background()
			user, err := svc.CreateUser(ctx, tc.input)

			if tc.wantErr {
				require.Error(t, err)
				require.Nil(t, user)
				if tc.expectedErr != nil {
					assert.ErrorIs(t, err, tc.expectedErr)
				}
			} else {
				require.NoError(t, err)
				require.NotNil(t, user)
				if tc.validate != nil {
					tc.validate(t, tc.input, user)
				}
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestLoginUser(t *testing.T) {
	const testUserID = "user-123"
	const testEmail = "user@example.com"
	const testPassword = "correctPassword123"

	existingUser := createTestUser(testUserID, testEmail, testPassword)

	testCases := []struct {
		name        string
		input       model.LoginUserParams
		setupMock   func(*MockUserRepository, *MockTokenMaker)
		wantErr     bool
		expectedErr error
		validate    func(t *testing.T, result *model.LoginUserResult)
	}{
		{
			name: "success - valid credentials",
			input: model.LoginUserParams{
				Email:    testEmail,
				Password: testPassword,
			},
			setupMock: func(mockRepo *MockUserRepository, mockToken *MockTokenMaker) {
				mockRepo.On("FindUserByEmail", mock.Anything, testEmail).
					Return(existingUser, nil)
				mockToken.On("CreateToken", testUserID).
					Return("jwt-token-xyz", nil)
				mockRepo.On("UpdateUserLastLogin", mock.Anything, testUserID, mock.AnythingOfType("time.Time")).
					Return(nil)
			},
			wantErr:     false,
			expectedErr: nil,
			validate: func(t *testing.T, result *model.LoginUserResult) {
				assert.Equal(t, "jwt-token-xyz", result.AccessToken)
			},
		},
		{
			name: "error - email not found",
			input: model.LoginUserParams{
				Email:    "nonexistent@example.com",
				Password: "anyPassword",
			},
			setupMock: func(mockRepo *MockUserRepository, mockToken *MockTokenMaker) {
				mockRepo.On("FindUserByEmail", mock.Anything, "nonexistent@example.com").
					Return(nil, model.ErrUserNotFound)
			},
			wantErr:     true,
			expectedErr: model.ErrIncorrectCredentials,
		},
		{
			name: "error - wrong password",
			input: model.LoginUserParams{
				Email:    testEmail,
				Password: "wrongPassword",
			},
			setupMock: func(mockRepo *MockUserRepository, mockToken *MockTokenMaker) {
				mockRepo.On("FindUserByEmail", mock.Anything, testEmail).
					Return(existingUser, nil)
			},
			wantErr:     true,
			expectedErr: model.ErrIncorrectCredentials,
		},
		{
			name: "error - token creation fails",
			input: model.LoginUserParams{
				Email:    testEmail,
				Password: testPassword,
			},
			setupMock: func(mockRepo *MockUserRepository, mockToken *MockTokenMaker) {
				mockRepo.On("FindUserByEmail", mock.Anything, testEmail).
					Return(existingUser, nil)
				mockToken.On("CreateToken", testUserID).
					Return("", errors.New("token creation failed"))
			},
			wantErr: true,
		},
		{
			name: "success - login succeeds even when updating last login fails",
			input: model.LoginUserParams{
				Email:    testEmail,
				Password: testPassword,
			},
			setupMock: func(mockRepo *MockUserRepository, mockToken *MockTokenMaker) {
				mockRepo.On("FindUserByEmail", mock.Anything, testEmail).
					Return(existingUser, nil)
				mockToken.On("CreateToken", testUserID).
					Return("jwt-token-xyz", nil)
				mockRepo.On("UpdateUserLastLogin", mock.Anything, testUserID, mock.AnythingOfType("time.Time")).
					Return(errors.New("update last login failed"))
			},
			wantErr:     false,
			expectedErr: nil,
			validate: func(t *testing.T, result *model.LoginUserResult) {
				assert.Equal(t, "jwt-token-xyz", result.AccessToken)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockRepo, mockToken, svc := newMocksAndService()
			tc.setupMock(mockRepo, mockToken)

			ctx := context.Background()
			result, err := svc.LoginUser(ctx, tc.input)

			if tc.wantErr {
				require.Error(t, err)
				require.Nil(t, result)
				if tc.expectedErr != nil {
					assert.ErrorIs(t, err, tc.expectedErr)
				}
			} else {
				require.NoError(t, err)
				require.NotNil(t, result)
				if tc.validate != nil {
					tc.validate(t, result)
				}
			}

			mockRepo.AssertExpectations(t)
			mockToken.AssertExpectations(t)
		})
	}
}
