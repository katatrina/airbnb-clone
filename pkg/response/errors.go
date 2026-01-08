package response

type ErrorCode string

const (
	CodeSuccess ErrorCode = "OK"

	// Client errors - Validation & Format (400)
	CodeValidationFailed  ErrorCode = "VALIDATION_FAILED"
	CodeInvalidJSONFormat ErrorCode = "INVALID_JSON_FORMAT"

	// Client errors - Authentication (401)
	CodeUnauthorized         ErrorCode = "UNAUTHORIZED"
	CodeIncorrectCredentials ErrorCode = "INCORRECT_CREDENTIALS"

	// Client errors - Authorization (403)
	CodeForbidden ErrorCode = "FORBIDDEN"

	// Client errors - Not Found (404)
	CodeNotFound ErrorCode = "NOT_FOUND"

	// Client errors - Conflict (409)
	CodeEmailAlreadyExists ErrorCode = "EMAIL_ALREADY_EXISTS"

	// Server errors (500)
	CodeInternalServerError ErrorCode = "INTERNAL_SERVER_ERROR"
)
