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
	"testing"

	"github.com/katatrina/airbnb-clone/services/user/internal/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

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
