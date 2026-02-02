// mock_test.go
// =============================================================================
// File này chứa các MOCK OBJECTS - những "đồ giả" thay thế dependencies thật.
//
// Tại sao cần mock?
// -----------------
// Khi test UserService, ta KHÔNG muốn:
//   - Kết nối database thật (chậm, cần setup)
//   - Tạo JWT thật (phức tạp, phụ thuộc config)
//
// Thay vào đó, ta tao "đồ giả" mà ta có thể KIỂM SOÁT hoàn toàn:
//   - Muốn nó trả về error? Được!
//   - Muốn nó trả về user cụ thể? Được!
//   - Muốn verify nó được gọi với params đúng? Được luôn!
//
// Thư viện testify/mock giúp ta làm điều này dễ dàng.
// ==========================================================================

package service

import (
	"context"
	"time"

	"github.com/katatrina/airbnb-clone/pkg/token"
	"github.com/katatrina/airbnb-clone/services/user/internal/model"
	"github.com/stretchr/testify/mock"
)

// MockUserRepository là bản giả của UserRepository.
// Nó "giả vờ" implement interface UserRepository, nhưng thay vì query DB thật,
// nó trả về những gì ta bảo nó trả về.
//
// mock.Mock là struct từ testify, nó có các methods:
//   - On(methodName, args...) : Setup expectation = "khi gọi method X với args Y"
//   - Return(values...)	   : "thì trả về values Z"
//   - Called(args)			   : Ghi nhận method được gọi, trả về values đã setup
type MockUserRepository struct {
	mock.Mock
}

// CheckEmailExists giả lập việc check email tồn tại trong DB.
//
// Giải thích flow:
//
//  1. Test code gọi: mockRepo.On("CheckEmailExists" ctx, "test@example.com").Return(true, nil)
//     -> Đăng ký expectation: "Nếu chỗ nào gọi CheckEmailExists với chính xác context và email này, trả về true và nil error"
//
//  2. Trong service method, code gọi: s.userRepo.CheckEmailExists(ctx, email)
//     -> Chạy vào mock function này
//
//  3. m.Called(ctx email) sẽ:
//     - Tìm expectation match với (ctx, email)
//     - Trả về giá trị đã đăng ký ở bước 1
//     - Ghi nhận là method này đã được gọi (để verify sau)
//     - Tuy hơi rườm rà nhưng đây là quy tắc của thư viện testify
func (m *MockUserRepository) CheckEmailExists(ctx context.Context, email string) (bool, error) {
	// m.Called() trả về mock.Arguments - một slice chứa các giá trị return đã được setup từ trước
	args := m.Called(ctx, email)

	// args.Bool(0) lấy giá trị thứ 0, cast sang bool
	// args.Error(1) lấy giá trị thứ 1, cast sang error
	return args.Bool(0), args.Error(1)
}

// CreateUser giả lập việc tạo user trong DB.
func (m *MockUserRepository) CreateUser(ctx context.Context, user model.User) (*model.User, error) {
	args := m.Called(ctx, user)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*model.User), args.Error(1)
}

// FindUserByEmail giả lập việc tìm user theo email.
//
// Lưu ý: Return value có thể là nil (khi user không tồn tại).
// Ta phải check nil trước khi cast, nếu không sẽ panic.
func (m *MockUserRepository) FindUserByEmail(ctx context.Context, email string) (*model.User, error) {
	args := m.Called(ctx, email)

	// args.Get(0) trả về interface{}, ta cần type assertion
	// Nhưng nếu ta setup Return(nil, err), thì args.Get(0) sẽ là nil
	// Và nil.(*model.User) sẽ panic!
	// Nên phải check nil trước:
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	// Giờ an toàn để cast
	return args.Get(0).(*model.User), args.Error(1)
}

// FindUserByID giả lập việc tìm user theo ID.
func (m *MockUserRepository) FindUserByID(ctx context.Context, id string) (*model.User, error) {
	args := m.Called(ctx, id)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*model.User), args.Error(1)
}

// UpdateUserLastLogin giả lập vệc update thời gian login cuối.
func (m *MockUserRepository) UpdateUserLastLogin(ctx context.Context, id string, lastLoginAt time.Time) error {
	// QUAN TRỌNG: Phải truyền TẤT CẢ parameters vào Called()
	// Nêếu thiếu, mock sẽ không match được expectation đã setup
	args := m.Called(ctx, id, lastLoginAt)

	return args.Error(0)
}

// MockTokenMaker giả lập pkg/token.TokenMaker interface.
// Thay vì tạo JWT thật (cần secret key, expiry config...),
// ta có thể bảo nó trả về bất kì token string nào ta muốn.
type MockTokenMaker struct {
	mock.Mock
}

// CreateToken giả lập việc tạo JWT token.
func (m *MockTokenMaker) CreateToken(userID string) (string, time.Time, error) {
	args := m.Called(userID)

	return args.String(0), args.Get(1).(time.Time), args.Error(2)
}

// VerifyToken giả lập việc verify JWT token.
// Hiện tại UserService không dùng method này, nhưng ta vẫn implement
// để satisfy interface TokenMaker.
func (m *MockTokenMaker) VerifyToken(tokenStr string) (*token.Claims, error) {
	args := m.Called(tokenStr)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*token.Claims), args.Error(1)
}
