package response

type ErrorCode string

const (
	CodeSuccess ErrorCode = "OK"

	// Validation & Format (400)
	CodeValidationFailed  ErrorCode = "VALIDATION_FAILED" // All input validation errors (body, URL, query)
	CodeJSONFormatInvalid ErrorCode = "INVALID_JSON_FORMAT"
	CodeReferenceInvalid  ErrorCode = "INVALID_REFERENCE" // Foreign key, relationship

	// Authentication (401)
	CodeAuthenticationRequired ErrorCode = "AUTHENTICATION_REQUIRED"
	CodeCredentialsInvalid     ErrorCode = "INVALID_CREDENTIALS"
	CodeTokenExpired           ErrorCode = "TOKEN_EXPIRED"
	CodeTokenInvalid           ErrorCode = "TOKEN_INVALID"

	// Authorization (403)

	// Not Found (404)
	CodeResourceNotFound ErrorCode = "RESOURCE_NOT_FOUND"

	// Conflict (409)
	CodeEmailAlreadyExists ErrorCode = "EMAIL_ALREADY_EXISTS"

	// Rate Limiting (429)
	CodeTooManyRequests ErrorCode = "TOO_MANY_REQUESTS"

	// Server errors (500)
	CodeInternalServerError ErrorCode = "INTERNAL_SERVER_ERROR"
	CodeServiceUnavailable  ErrorCode = "SERVICE_UNAVAILABLE" // DB down, third-party API down
)
