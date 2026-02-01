package model

import "time"

type CreateUserParams struct {
	DisplayName string
	Email       string
	Password    string
}

type LoginUserParams struct {
	Email    string
	Password string
}

type LoginUserResult struct {
	AccessToken          string
	AccessTokenExpiresAt time.Time
	User                 *User
}
