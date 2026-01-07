// Package handler contains HTTP handlers and DTOs.
// DTOs (Data Transfer Objects) define the shape of HTTP request/response bodies.
//
// Why separate DTOs for Handler vs Service?
// - Handler DTOs have `json` tags for serialization and `binding` tags for validation
// - Service DTOs are pure Go structs without any HTTP concerns
// - This separation allows the service to be reused with gRPC, CLI, etc.
package handler

type RegisterRequest struct {
	DisplayName string `json:"displayName" binding:"required,min=2,max=100,safename"`
	Email       string `json:"email" binding:"required,email,max=255"`
	Password    string `json:"password" binding:"required,min=8,max=72,strongpass"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
	AccessToken string `json:"accessToken"`
}

// UserResponse is the standard user representation in API responses.
// Used by GET /users/me and other endpoints that return user data.
type UserResponse struct {
	ID          string `json:"id"`
	DisplayName string `json:"displayName"`
	Email       string `json:"email"`
	CreatedAt   int64  `json:"createdAt"`
}
