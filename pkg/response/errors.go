package response

type ErrorCode string

const (
	CodeSuccess ErrorCode = "OK"

	CodeValidationFailed  ErrorCode = "VALIDATION_FAILED" // All input validation errors (body, URL, query)
	CodeJSONFormatInvalid ErrorCode = "INVALID_JSON_FORMAT"
	CodeReferenceInvalid  ErrorCode = "INVALID_REFERENCE" // Foreign key, relationship

	CodeListingNotDraft              ErrorCode = "LISTING_NOT_DRAFT"
	CodeListingNotActive             ErrorCode = "LISTING_NOT_ACTIVE"
	CodeListingIncomplete            ErrorCode = "LISTING_INCOMPLETE"
	CodeActiveListingCannotBeUpdated ErrorCode = "ACTIVE_LISTING_CANNOT_BE_UPDATED"

	CodeAuthenticationRequired ErrorCode = "AUTHENTICATION_REQUIRED"
	CodeCredentialsInvalid     ErrorCode = "INVALID_CREDENTIALS"
	CodeTokenExpired           ErrorCode = "TOKEN_EXPIRED"
	CodeTokenInvalid           ErrorCode = "TOKEN_INVALID"

	CodeUserNotFound     ErrorCode = "USER_NOT_FOUND"
	CodeListingNotFound  ErrorCode = "LISTING_NOT_FOUND"
	CodeProvinceNotFound ErrorCode = "PROVINCE_NOT_FOUND"
	CodeDistrictNotFound ErrorCode = "DISTRICT_NOT_FOUND"
	CodeWardNotFound     ErrorCode = "WARD_NOT_FOUND"

	CodeEmailAlreadyExists ErrorCode = "EMAIL_ALREADY_EXISTS"

	CodeTooManyRequests ErrorCode = "TOO_MANY_REQUESTS"

	CodeInternalServerError ErrorCode = "INTERNAL_SERVER_ERROR"
	CodeServiceUnavailable  ErrorCode = "SERVICE_UNAVAILABLE"
)
