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
	"fmt"
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
	// Email must be unique
	existing, err := s.userRepo.FindUserByEmail(ctx, arg.Email)
	if err != nil && !errors.Is(err, model.ErrUserNotFound) {
		return nil, err
	}
	if existing != nil {
		return nil, model.ErrEmailAlreadyExists
	}

	// Password must be hashed (bcrypt cost 10)
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(arg.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	// User starts as unverified (implementing email verification)
	userID, _ := uuid.NewV7()
	now := time.Now()

	user := model.User{
		ID:            userID.String(),
		DisplayName:   arg.DisplayName,
		Email:         arg.Email,
		PasswordHash:  string(hashedPassword),
		EmailVerified: false,
		CreatedAt:     now,
		UpdatedAt:     now,
	}

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
		// Don't reveal if email exists or not (security)
		if errors.Is(err, model.ErrUserNotFound) {
			return nil, model.ErrIncorrectCredentials
		}
		return nil, err
	}

	// Check if account is active
	if user.DeletedAt != nil {
		return nil, model.ErrIncorrectCredentials
	}

	// Compare password with stored hash
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(arg.Password))
	if err != nil {
		// Password doesn't match - return same generic error
		// Don't reveal that the password is wrong (email was correct)
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return nil, model.ErrIncorrectCredentials
		}
		return nil, err
	}

	// Generate access token
	accessToken, err := s.tokenMaker.CreateToken(user.ID)
	if err != nil {
		return nil, err
	}

	// Update last login time
	now := time.Now()
	user.LastLoginAt = &now
	err = s.userRepo.UpdateLastLogin(ctx, user.ID, user.LastLoginAt)

	return &LoginResult{
		AccessToken: accessToken,
	}, nil
}

func (s *UserService) GetUserByID(ctx context.Context, id string) (*model.User, error) {
	return s.userRepo.FindUserByID(ctx, id)
}
