package token

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

const MinHMACSecretKeySize = 32

// JWTMaker is a TokenMaker implementation that uses JSON Web Tokens (JWT).
// It uses the HS256 (HMAC-SHA256) signing algorithm with a symmetric secret key.
//
// Why HS256?
// - Simple: One secret key for both signing and verification
// - Fast: Symmetric encryption is faster than asymmetric (RS256)
// - Sufficient: For monolithic or internal microservices where you control all services
//
// When to consider RS256 (asymmetric):
// - When third parties need to verify tokens but shouldn't be able to create them
// - When you want to rotate keys without coordinating across services
type JWTMaker struct {
	secretKey []byte
	expiry    time.Duration
}

var _ TokenMaker = (*JWTMaker)(nil)

// NewJWTMaker creates a new JWTMaker with the given secret key and expiry duration.
//
// Parameters:
//   - secretKey: The secret used to sign tokens. Should be at least 32 characters
//     (256 bits) for security. In production, load this from environment variables.
//   - expiry: How long tokens are valid. Common values: 15m for access tokens,
//     7d for refresh tokens.
//
// Example:
//
//	maker := token.NewJWTMaker(os.Getenv("JWT_SECRET"), 24*time.Hour)
func NewJWTMaker(secretKey []byte, expiry time.Duration) (*JWTMaker, error) {
	if len(secretKey) < MinHMACSecretKeySize {
		return nil, fmt.Errorf("secret key must be at least %d bytes", MinHMACSecretKeySize)
	}

	return &JWTMaker{
		secretKey: secretKey,
		expiry:    expiry,
	}, nil
}

// CreateToken generates a new JWT for the given user ID.
//
// The token contains standard JWT claims (RFC 7519):
//   - sub (Subject): The user ID - who this token represents
//   - exp (Expiration): When the token expires
//   - iat (Issued At): When the token was created
//   - nbf (Not Before): Token is not valid before this time (set to now)
//   - jti (JWT ID): Unique identifier for this token (for revocation lists)
//
// The token is signed using HS256 (HMAC-SHA256) algorithm.
func (m *JWTMaker) CreateToken(userID string) (string, time.Time, error) {
	now := time.Now()
	expiresAt := now.Add(m.expiry)

	// RegisteredClaims is from jwt/v5 and follows RFC 7519 standard claim names.
	// Using standard claims makes your tokens interoperable with other systems.
	claims := jwt.RegisteredClaims{
		// Subject identifies the principal (user) this token represents.
		// This is the most important claim - it's how you know who is authenticated.
		Subject: userID,

		// ExpiresAt is when this token becomes invalid.
		// After this time, VerifyToken will return ErrTokenExpired.
		ExpiresAt: jwt.NewNumericDate(expiresAt),

		// IssuedAt is when this token was created.
		// Useful for debugging and for implementing "logout everywhere" features
		// (invalidate all tokens issued before a certain time).
		IssuedAt: jwt.NewNumericDate(now),

		// NotBefore means the token is not valid before this time.
		// We set it to now, meaning the token is immediately valid.
		// You could set it to a future time for delayed activation.
		NotBefore: jwt.NewNumericDate(now),

		// ID is a unique identifier for this specific token.
		// Useful for:
		// - Token revocation (maintain a blacklist of revoked JTIs)
		// - Preventing replay attacks
		// - Debugging (trace a specific token in logs)
		ID: uuid.NewString(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenStr, err := token.SignedString(m.secretKey)
	if err != nil {
		return "", time.Time{}, err
	}

	return tokenStr, expiresAt, nil
}

// VerifyToken parses a JWT string and validates it.
//
// Validation includes:
//   - Signature verification (was this token created with our secret key?)
//   - Expiration check (has the token expired?)
//   - Not Before check (is the token active yet?)
//   - Algorithm check (prevent algorithm confusion attacks)
//
// Returns Claims if the token is valid, or an error if validation fails.
func (m *JWTMaker) VerifyToken(tokenString string) (*Claims, error) {
	// ParseWithClaims parses the token and validates the signature.
	// The callback function provides the key for signature verification.
	token, err := jwt.ParseWithClaims(
		tokenString,
		&jwt.RegisteredClaims{},
		func(token *jwt.Token) (any, error) {
			// This function is called during parsing to get the verification key.
			// You could also check token.Method here to ensure the expected algorithm.
			return []byte(m.secretKey), nil
		},
		// WithValidMethods explicitly specifies which signing algorithms are allowed.
		// This prevents "algorithm confusion" attacks where an attacker changes
		// the algorithm in the header to bypass signature verification.
		jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}),
	)

	// Handle parsing/validation errors
	if err != nil {
		// Check if the error is specifically about expiration.
		// jwt/v5 returns specific error types we can check.
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, ErrTokenExpired
		}
		// For all other errors (invalid signature, malformed, etc.)
		return nil, ErrTokenInvalid
	}

	// Extract the claims from the parsed token.
	// We already know the token is valid at this point.
	claims, ok := token.Claims.(*jwt.RegisteredClaims)
	if !ok {
		// This shouldn't happen if ParseWithClaims succeeded, but be defensive.
		return nil, ErrTokenInvalid
	}

	// Convert JWT claims to our application's Claims struct.
	// This decouples our application from the JWT library's types.
	return &Claims{
		UserID:    claims.Subject,
		IssuedAt:  claims.IssuedAt.Time,
		ExpiresAt: claims.ExpiresAt.Time,
	}, nil
}
