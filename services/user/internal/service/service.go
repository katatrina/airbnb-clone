// Package service contains business logic.
// This layer is responsible for:
// 1. Orchestrating calls to repository
// 2. Implementing business rules (password hashing, etc.)
// 3. Transforming between DTOs and domain models
//
// Key principle: Service layer knows NOTHING about HTTP.
// No gin.Context, no HTTP status codes, no JSON tags.
// This makes it reusable across different transports (HTTP, gRPC, CLI, etc.)
package service

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/katatrina/airbnb-clone/pkg/token"
	"github.com/katatrina/airbnb-clone/services/user/internal/model"
	"github.com/katatrina/airbnb-clone/services/user/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	userRepo   *repository.UserRepository
	tokenMaker token.TokenMaker
}

func NewUserService(userRepo *repository.UserRepository, tokenMaker token.TokenMaker) *UserService {
	return &UserService{
		userRepo:   userRepo,
		tokenMaker: tokenMaker,
	}
}

type CreateUserParams struct {
	DisplayName string
	Email       string
	Password    string
}

func (s *UserService) CreateUser(ctx context.Context, arg CreateUserParams) (*model.User, error) {
	userID, err := uuid.NewV7()
	if err != nil {
		return nil, err
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(arg.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	now := time.Now()
	user := model.User{
		ID:           userID.String(),
		DisplayName:  arg.DisplayName,
		Email:        arg.Email,
		PasswordHash: string(hashedPassword),
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	// Repository returns model.ErrEmailAlreadyExists if email is taken
	// We pass this through to the handler - no need to wrap or transform
	err = s.userRepo.CreateUser(ctx, &user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

type LoginParams struct {
	Email    string
	Password string
}

type LoginResult struct {
	AccessToken string
}

func (s *UserService) Login(ctx context.Context, arg LoginParams) (*LoginResult, error) {
	// Find user by email
	user, err := s.userRepo.FindUserByEmail(ctx, arg.Email)
	if err != nil {
		// If user not found, return generic "invalid credentials"
		// Don't reveal that the email doesn't exist
		if errors.Is(err, model.ErrUserNotFound) {
			return nil, model.ErrInvalidCredentials
		}
		return nil, err
	}

	// Compare password with stored hash
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(arg.Password))
	if err != nil {
		// Password doesn't match - return same generic error
		// Don't reveal that the password is wrong (email was correct)
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return nil, model.ErrInvalidCredentials
		}
		return nil, err
	}

	// Generate access token
	accessToken, err := s.tokenMaker.CreateToken(user.ID)
	if err != nil {
		return nil, err
	}

	return &LoginResult{
		AccessToken: accessToken,
	}, nil
}

// GetUserByID retrieves a user by their ID.
// Pass through model.ErrUserNotFound if user doesn't exist.
func (s *UserService) GetUserByID(ctx context.Context, id string) (*model.User, error) {
	return s.userRepo.FindUserByID(ctx, id)
}
