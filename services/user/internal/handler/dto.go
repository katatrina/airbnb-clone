// Package handler contains HTTP handlers and DTOs.
// DTOs (Data Transfer Objects) define the shape of HTTP request/response bodies.
//
// Why separate DTOs for Handler vs Service?
// - Handler DTOs have `json` tags for serialization and `binding` tags for validation
// - Service DTOs are pure Go structs without any HTTP concerns
// - This separation allows the service to be reused with gRPC, CLI, etc.
package handler

type RegisterRequest struct {
	Email string `json:"email" binding:"required,email,max=255" normalize:"trim,lower"`
	// Password: bcrypt only uses first 72 bytes, so we limit max length to 72
	Password    string `json:"password" binding:"required,min=8,maxbytes=72,strongpass" normalize:"trim"`
	DisplayName string `json:"displayName" binding:"required,min=2,max=100,displayname" normalize:"trim"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email" normalize:"trim,lower"`
	Password string `json:"password" binding:"required" normalize:"trim"`
	// Password: NO min/max validation on login (accept any length)
	// User might have old password before we added min=8 rule
}

type LoginResponse struct {
	AccessToken string `json:"accessToken"`
}

// UserResponse is the standard user representation in API responses.
// Used by GET /users/me and other endpoints that return user data.
type UserResponse struct {
	ID            string `json:"id"`
	DisplayName   string `json:"displayName"`
	Email         string `json:"email"`
	EmailVerified bool   `json:"emailVerified"`
	LastLoginAt   *int64 `json:"lastLoginAt,omitempty"`
	CreatedAt     int64  `json:"createdAt"`
}
