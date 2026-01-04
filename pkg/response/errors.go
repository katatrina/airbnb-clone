package response

// ============================================================================
// Error Codes
// ============================================================================
// These codes help clients handle errors programmatically.
// They follow a pattern: first digit = HTTP status category, rest = specific error.
//
// Pattern:
// - 4xxx = Client errors (4xx HTTP status)
// - 5xxx = Server errors (5xx HTTP status)
//
// This allows the frontend to:
// 1. Check success/failure via the `success` field
// 2. Handle specific errors via the `code` field

const (
	CodeSuccess = 0 // ✅

	CodeValidationFailed = 4001
	CodeInvalidJSON      = 4002
	CodeMissingField     = 4003
	CodeInvalidFormat    = 4004

	CodeUnauthorized = 4010 // ✅
	CodeTokenExpired = 4011
	CodeTokenInvalid = 4012

	CodeForbidden        = 4030
	CodeInsufficientRole = 4031

	CodeNotFound     = 4040
	CodeUserNotFound = 4041

	CodeConflict    = 4090
	CodeEmailExists = 4091 // ✅

	CodeInternalError   = 5000 // ✅
	CodeDatabaseError   = 5001
	CodeExternalService = 5002
)

//var ErrorMessages = map[int]string{
//	CodeSuccess:          "OK",
//	CodeValidationFailed: "Validation failed",
//	CodeEmailExists:      "Email already exists",
//	CodeUserNotFound:     "User not found",
//}
