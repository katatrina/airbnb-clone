package response

const (
	CodeSuccess = 0

	CodeValidationFailed = 4001
	CodeInvalidJSON      = 4002
	CodeMissingField     = 4003
	CodeInvalidFormat    = 4004

	CodeUnauthorized = 4010
	CodeTokenExpired = 4011
	CodeTokenInvalid = 4012

	CodeForbidden        = 4030
	CodeInsufficientRole = 4031

	CodeNotFound         = 4040
	CodeUserNotFound     = 4041
	CodeListingNotFound  = 4042
	CodeProvinceNotFound = 4043

	CodeConflict      = 4090
	CodeEmailExists   = 4091
	CodeListingExists = 4092

	CodeInternalError   = 5000
	CodeDatabaseError   = 5001
	CodeExternalService = 5002
)

var ErrorMessages = map[int]string{
	CodeSuccess:          "OK",
	CodeValidationFailed: "Validation failed",
	CodeEmailExists:      "Email already exists",
	CodeUserNotFound:     "User not found",
}
