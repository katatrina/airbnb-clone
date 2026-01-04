// Package constant contains shared constants used across the application.
// Using constants instead of magic strings prevents typos and makes refactoring easier.
package constant

// UserIDKey is the context key for storing/retrieving the authenticated user's ID.
// This is set by AuthMiddleware and read by handlers that need the current user.
const UserIDKey = "userID"
