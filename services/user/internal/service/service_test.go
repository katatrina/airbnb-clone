// service_test.go
// =============================================================================
// Unit Tests cho UserService
//
// File này test 3 methods của UserService:
//  1. CreateUser  - Đăng ký user mới
//  2. LoginUser   - Đăng nhập
//  3. GetUserByID - Lấy thông tin user
//
// Pattern sử dụng:
//   - Table-Driven Tests: Gom nhiều test cases vào một bảng, loop qua
//   - AAA Pattern: Arrange -> Act -> Assert
//
// Chạy tests:
//
//	go test -v -cover ./internal/service
//
// =============================================================================
package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/katatrina/airbnb-clone/services/user/internal/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)

// createTestUser tạo một *model.User với các giá trị mặc định hợp lệ.
// Bạn có thể override bất kì field nào sau khi gọi function này.
//
// Tại sao cần helper này?
// Vì mỗi test case đều cần một User object, và việc tạo thủ công mỗi lần
// rất dài dòng và dễ quên field. Helper này đảm bảo ta luôn có
// một User "đầy đủ" để test.
func createTestUser(id, email, password string) *model.User {
	// Hash password giống như production code làm
	// Điều này QUAN TRỌNG vì LoginUser sẽ dùng bcrypt.CompareHashAndPassword
	// Nếu password không được hash đúng cách, test sẽ fail
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	now := time.Now()

	return &model.User{
		ID:            id,
		DisplayName:   "Test User",
		Email:         email,
		PasswordHash:  string(hashedPassword),
		EmailVerified: false,
		LastLoginAt:   nil, // Chưa login lần nào
		CreatedAt:     now,
		UpdatedAt:     now,
	}
}

// newMocksAndService tạo mock objects và UserService instance.
// Đây là "factory function" giúp setup test nhanh hơn.
//
// Returns:
//   - *MockUserRepository: Mock để setup expectations cho repository
//   - *MockTokenMaker: Mock để setup expectations cho token
//   - *UserService: Service instance với mocks đã inject
func newMocksAndService() (*MockUserRepository, *MockTokenMaker, *UserService) {
	mockRepo := new(MockUserRepository)
	mockToken := new(MockTokenMaker)

	// Inject mocks vào service thông qua constructor
	// Đây là Dependency Injection pattern - service không biết
	// nó đang dùng mock hay real implementation
	service := NewUserService(mockRepo, mockToken)

	return mockRepo, mockToken, service
}

// TestCreateUser .
//
// CreateUser có scenarios cần test:
//  1. Success - Happy path, mọi thứ OK
//  2. Email already exists - CheckEmailExists trả về true
func TestCreateUser(t *testing.T) {
	// Định nghĩa test cases bằng table-driven pattern
	// Mỗi test case là một struct chứa:
	//	- name: Tên test case (hiển thị khi chạy test)
	//	- input: Dữ liệu đầu vào
	//	- setupMock: Function để setup mock expectations
	//	- wantErr: Có expect error không?
	//	- expectedErr: Error cụ tể expect nhận được
	//	- validate: Function để validate kết quả (optional)
	testCases := []struct {
		// name mô tả ngắn gọn test case đang test gì
		// Convention: "should <expected bahavior> when <condition>"
		name string

		// input là CreateUserParams - đầu vào cho CreateUser
		input model.CreateUserParams

		// setupMock là function nhận mock objects và setup expectations
		// Tại sao dùng function thay vì data?
		// Vì mỗi test case cần setup KHÁC NHAU, và setup có thể phức tạp
		// (multiple calls, different return values, etc.)
		setupMock func(*MockUserRepository, *MockTokenMaker)

		// wantErr cho biết ta EXPECT test case này có error hay không
		// true = expect error, false = expect success
		wantErr bool

		// expectedErr là error cụ thể ta expect
		// Dùng với errors.Is() để check error type
		expectedErr error

		// validate là function để kiểm tra kết quả chi tiết hơn
		// Chỉ chạy khi wantErr = false (success case)
		// Nhận user trả về và input để so sánh
		validate func(t *testing.T, input model.CreateUserParams, user *model.User)
	}{ // Test case 1: Happy path - Mọi thứ OK
		{
			name: "success - creates user with valid input",
			input: model.CreateUserParams{
				DisplayName: "New User",
				Email:       "newuser@example.com",
				Password:    "strongpassword123",
			},
			setupMock: func(mockRepo *MockUserRepository, mockToken *MockTokenMaker) {
				// Expect: CheckEmailExists được gọi với email này
				// Return: false (email chưa tồn tại), nil (không có error)
				mockRepo.On("CheckEmailExists", mock.Anything, "newuser@example.com").
					Return(false, nil)

				// Expect: CreateUser được gọi với BẤT KỲ *model.User nào
				// Tai sao dùng mock.AnythingOfType thay vì pass vào user cụ thể?
				// Vì service tạo user với ID random (uuid.NewV7), ta không biết trước
				// Ta chỉ cần verify là CreateUser ĐƯỢC GỌI, không cần match exact value
				mockRepo.On("CreateUser", mock.Anything, mock.AnythingOfType("*model.User")).
					Return(nil) // Success
			},
			wantErr:     false,
			expectedErr: nil,
			validate: func(t *testing.T, input model.CreateUserParams, user *model.User) {
				// verify các field được set đúng
				assert.Equal(t, input.Email, user.Email)
				assert.Equal(t, input.DisplayName, user.DisplayName)

				// ID phải được generate (không empty)
				assert.NotEmpty(t, user.ID)

				// Password phải được HASH, không phải plaintext!
				// Đây là security test quan trọng
				assert.NotEqual(t, input.Password, user.PasswordHash)

				// bcrypt hash luôn bắt đầu bằng "$2a$" hoặc "$2b$"
				assert.Contains(t, user.PasswordHash, "$2a$")

				// User mới chưa verify email
				assert.False(t, user.EmailVerified)

				// Timestamp phải được set
				assert.False(t, user.CreatedAt.IsZero())
				assert.False(t, user.UpdatedAt.IsZero())
			},
		},

		// Test case 2: Email đã tồn tại
		{
			name: "error - email already exists",
			input: model.CreateUserParams{
				DisplayName: "Existing User",
				Email:       "existing@example.com",
				Password:    "strongpassword123",
			},
			setupMock: func(mockRepo *MockUserRepository, mockToken *MockTokenMaker) {
				// CheckEmailExists trả về TRUE = email đã tồn tại
				mockRepo.On("CheckEmailExists", mock.Anything, "existing@example.com").
					Return(true, nil)

				// KHÔNG setup CreateUser vì nó không được gọi
				// Service phải return error TRƯỚC KHI gọi CreateUser
			},
			wantErr:     true,
			expectedErr: model.ErrEmailAlreadyExists,
			validate:    nil, // Không cần validate vì expect error
		},
	}

	// Đã chuẩn bị xong tất cả các test case
	// Giờ chỉ việc chạy thôi
	for _, tc := range testCases {
		// t.Run tạo một "sub-test" cho mỗi test case
		// Khi chạy `go test -v`, bạn sẽ thấy kiểu như này:
		//   === RUN   TestCreateUser/success_-_creates_user_with_valid_input
		//   === RUN   TestCreateUser/error_-_email_already_exists
		//   ...
		t.Run(tc.name, func(t *testing.T) {
			// Arrange
			// Setup mocks và service cho test case này
			mockRepo, mockToken, svc := newMocksAndService()

			// Gọi setupMock để đăng ký expectations
			tc.setupMock(mockRepo, mockToken)

			// Act
			// Gọi function cần test
			ctx := context.Background()
			user, err := svc.CreateUser(ctx, tc.input)

			// Assert
			if tc.wantErr {
				// Expect error
				require.Error(t, err, "Expected error but got nil")
				require.Nil(t, user, "User should be nil when error occurs")

				// Nếu có expectedErr cụ thể, check với errors.Is
				if tc.expectedErr != nil {
					assert.ErrorIs(t, err, tc.expectedErr)
				}
			} else {
				// Expect success
				require.NoError(t, err, "Expected no error but got: %v", err)
				require.NotNil(t, user, "User should not be nil on success")

				// Chạy custom validation nếu có
				if tc.validate != nil {
					tc.validate(t, tc.input, user)
				}
			}

			// Verify rằng tất cả expectations đã được gọi đúng
			// Nếu ta setup mock.On("CreateUser",...) mà service không gọi CreateUser,
			// AssertExpectations sẽ fail
			mockRepo.AssertExpectations(t)
		})
	}
}

// TestLoginUser .
// LoginUser có các scenarios cần test:
//  1. Success - Login thành công, trả về token
//  2. User not found - Email không tồn tại
//  3. Wrong password - Password sai
//  4. Token creation fails - TokenMaker.CreateToken fails
//  5. UpdateUserLastLogin fails - Vẫn phải trả về success (non-critical)
func TestLoginUser(t *testing.T) {
	// Tạo password và user MỘT LẦN, dùng chung cho nhiều test cases
	// Đây là optimization: bcrypt.GenerateFromPassword chậm, không nên chạy nhiều lần
	const testUserID = "user-123"
	const testEmail = "user@example.com"
	const testPassword = "correctPassword123"

	// QUAN TRỌNG: User phải có password hash THẬT từ bcrypt
	// Vì service code sẽ gọi bcrypt.CompareHashAndPassword
	// Lưu ý: Mặc dù hai hash password khác giá trị do bcrypt random salt,
	// nhưng khi so sánh chúng đều là một password.
	existingUser := createTestUser(testUserID, testEmail, testPassword)

	testCases := []struct {
		name        string
		input       model.LoginUserParams
		setupMock   func(*MockUserRepository, *MockTokenMaker)
		wantErr     bool
		expectedErr error
		validate    func(t *testing.T, result *model.LoginUserResult)
	}{ // Test case 1: Happy path - Login thành công
		{
			name: "success - valid credentials",
			input: model.LoginUserParams{
				Email:    testEmail,
				Password: testPassword, // Password ĐÚNG
			},
			setupMock: func(mockRepo *MockUserRepository, mockToken *MockTokenMaker) {
				// FindUserByEmail trả về user với password hash
				mockRepo.On("FindUserByEmail", mock.Anything, testEmail).
					Return(existingUser, nil)

				// TokenMaker.CreateToken được gọi với user ID
				mockToken.On("CreateToken", testUserID).
					Return("jwt-token-xyz", nil)

				// UpdateUserLastLogin được gọi (non-critical nên dùng mock.Anything cho time)
				mockRepo.On("UpdateUserLastLogin", mock.Anything, testUserID, mock.AnythingOfType("time.Time")).
					Return(nil)
			},
			wantErr:     false,
			expectedErr: nil,
			validate: func(t *testing.T, result *model.LoginUserResult) {
				assert.Equal(t, "jwt-token-xyz", result.AccessToken)
			},
		},

		// Test case 2: Email không tồn tại
		{
			name: "error  - email not found",
			input: model.LoginUserParams{
				Email:    "nonexistent@example.com",
				Password: "anypassword",
			},
			setupMock: func(mockRepo *MockUserRepository, mockToken *MockTokenMaker) {
				// FindUserByEmail trả về ErrUserNotFound
				mockRepo.On("FindUserByEmail", mock.Anything, "nonexistent@example.com").
					Return(nil, model.ErrUserNotFound)

				// Không cần setup CreateToken vì không được gọi
			},
			wantErr: true,
			// QUAN TRỌNG: Service trả về ErrIncorrectCredentials, KHÔNG phải ErrUserNotFound
			// Đây là security best practices: không cho attacker biết email có tồn tại hay không
			expectedErr: model.ErrIncorrectCredentials,
			validate:    nil, // không cần validate khi có error
		},

		// Test case 3: Password sai
		{
			name: "error - wrong password",
			input: model.LoginUserParams{
				Email:    testEmail,
				Password: "wrongPassword",
			},
			setupMock: func(mockRepo *MockUserRepository, mockToken *MockTokenMaker) {
				// FindUserByEmail trả về user với password hash ĐÚNG
				mockRepo.On("FindUserByEmail", mock.Anything, testEmail).
					Return(existingUser, nil)

				// CreateToken không được gọi vì password check fail trước
			},
			wantErr:     true,
			expectedErr: model.ErrIncorrectCredentials,
			validate:    nil,
		},

		// Test case 4: Token creation fails
		{
			name: "error - token creation fails",
			input: model.LoginUserParams{
				Email:    testEmail,
				Password: testPassword,
			},
			setupMock: func(mockRepo *MockUserRepository, mockToken *MockTokenMaker) {
				mockRepo.On("FindUserByEmail", mock.Anything, testEmail).
					Return(existingUser, nil)

				// TokenMaker trả về error
				mockToken.On("CreateToken", testUserID).
					Return("", errors.New("random error when creating token I didn't know about"))
			},
			wantErr:     true,
			expectedErr: nil,
			validate:    nil,
		},

		// Test case 6: UpdateUserLastLogin fails nhưng login vẫn thành công
		// ---------------------------------------------------------------------
		// Đây là test case QUAN TRỌNG về business logic:
		// - Việc update last login là "nice to have", không critical
		// - Nếu nó fail, user vẫn phải login được
		// - Service chỉ log warning, không return error
		// ---------------------------------------------------------------------
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

				// UpdateUserLastLogin fails
				mockRepo.On("UpdateUserLastLogin", mock.Anything, testUserID, mock.AnythingOfType("time.Time")).
					Return(errors.New("random error when updating user's last login timestamp"))
			},
			wantErr:     false, // VẪN SUCCESS!
			expectedErr: nil,
			validate: func(t *testing.T, result *model.LoginUserResult) {
				// Token vẫn được trả về dù UpdateUserLastLogin fail
				assert.Equal(t, "jwt-token-xyz", result.AccessToken)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Arrange
			mockRepo, mockToken, svc := newMocksAndService()
			tc.setupMock(mockRepo, mockToken)

			// Act
			ctx := context.Background()
			result, err := svc.LoginUser(ctx, tc.input)

			// Assert
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
