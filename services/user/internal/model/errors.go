package model

import "errors"

var (
	ErrIncorrectCredentials = errors.New("incorrect email or password")

	ErrUserNotFound = errors.New("user not found")

	ErrEmailAlreadyExists = errors.New("email already exists")
)
