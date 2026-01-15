package handler

type RegisterRequest struct {
	// Lowercase to prevent duplicates
	Email string `json:"email" validate:"required,email,max=255" normalize:"trim,lower"`

	// Min 8 chars, max 72 bytes (bcrypt limit)
	Password string `json:"password" validate:"required,min=8,maxbytes=72"`

	// Public display name
	DisplayName string `json:"displayName" validate:"required,min=2,max=100" normalize:"trim,singlespace"`
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email" normalize:"trim,lower"`
	Password string `json:"password" validate:"required"`
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
