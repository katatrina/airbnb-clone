package service

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/katatrina/airbnb-clone/pkg/token"
	"github.com/katatrina/airbnb-clone/services/user/internal/model"
	"golang.org/x/crypto/bcrypt"
)

type UserRepository interface {
	CreateUser(ctx context.Context, user *model.User) error
	FindUserByEmail(ctx context.Context, email string) (*model.User, error)
	FindUserByID(ctx context.Context, id string) (*model.User, error)
	UpdateUserLastLogin(ctx context.Context, id string, lastLoginAt time.Time) error
	CheckEmailExists(ctx context.Context, email string) (bool, error)
}

type UserService struct {
	userRepo   UserRepository
	tokenMaker token.TokenMaker
}

func NewUserService(userRepo UserRepository, tokenMaker token.TokenMaker) *UserService {
	return &UserService{
		userRepo:   userRepo,
		tokenMaker: tokenMaker,
	}
}

func (s *UserService) CreateUser(ctx context.Context, arg model.CreateUserParams) (*model.User, error) {
	exists, err := s.userRepo.CheckEmailExists(ctx, arg.Email)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, model.ErrEmailAlreadyExists
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(arg.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

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

func (s *UserService) LoginUser(ctx context.Context, arg model.LoginUserParams) (*model.LoginUserResult, error) {
	user, err := s.userRepo.FindUserByEmail(ctx, arg.Email)
	if err != nil {
		if errors.Is(err, model.ErrUserNotFound) {
			return nil, model.ErrIncorrectCredentials
		}
		return nil, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(arg.Password))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return nil, model.ErrIncorrectCredentials
		}
		return nil, err
	}

	accessToken, err := s.tokenMaker.CreateToken(user.ID)
	if err != nil {
		return nil, err
	}

	now := time.Now()
	err = s.userRepo.UpdateUserLastLogin(ctx, user.ID, now)
	if err != nil {
		log.Printf("[WARN] Failed to update last login for user: %v", err)
	}

	return &model.LoginUserResult{
		AccessToken: accessToken,
	}, nil
}

func (s *UserService) GetUserByID(ctx context.Context, id string) (*model.User, error) {
	return s.userRepo.FindUserByID(ctx, id)
}
