package handler

type RegisterRequest struct {
<<<<<<< HEAD
	// Lowercase to prevent duplicates
	Email string `json:"email" binding:"required,email,max=255" normalize:"trim,lower"`

	// Min 8 chars, max 72 bytes (bcrypt limit)
	Password string `json:"password" binding:"required,min=8,maxbytes=72"`

	// Public display name
	DisplayName string `json:"displayName" binding:"required,min=2,max=100" normalize:"trim,singlespace"`
=======
	Email string `json:"email" binding:"required,email,max=255" normalize:"trim,lower"`
	// Password: bcrypt only uses first 72 bytes, so we limit max length to 72
	Password    string `json:"password" binding:"required,min=8,maxbytes=72,strongpass" normalize:"trim"`
	DisplayName string `json:"displayName" binding:"required,min=2,max=100,displayname" normalize:"trim"`
>>>>>>> dfabe7596391d9c6c7bf9d1e24a4534522056979
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email" normalize:"trim,lower"`
<<<<<<< HEAD
	Password string `json:"password" binding:"required"`
=======
	Password string `json:"password" binding:"required" normalize:"trim"`
>>>>>>> dfabe7596391d9c6c7bf9d1e24a4534522056979
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
