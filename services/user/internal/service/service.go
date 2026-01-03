package service

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/katatrina/airbnb-clone/services/user/internal/model"
	"github.com/katatrina/airbnb-clone/services/user/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	userRepo *repository.UserRepository
}

func NewUserService(userRepo *repository.UserRepository) *UserService {
	return &UserService{userRepo: userRepo}
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

	err = s.userRepo.CreateUser(ctx, &user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}
