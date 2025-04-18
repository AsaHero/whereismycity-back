package outerr

// API Error Codes
const (
	// Common errors
	CodeBadRequest      = "BAD_REQUEST"
	CodeUnauthorized    = "UNAUTHORIZED"
	CodeForbidden       = "FORBIDDEN"
	CodeNotFound        = "NOT_FOUND"
	CodeConflict        = "CONFLICT"
	CodeInternalError   = "INTERNAL_ERROR"
	CodeNoChanges       = "NOT_MODIFIED"
	CodeTooManyRequests = "TOO_MANY_REQUESTS"

	// Validation errors
	CodeValidation        = "VALIDATION_ERROR"
	CodeInvalidValue      = "INVALID_VALUE"
	CodeInvalidFormat     = "INVALID_FORMAT"
	CodeInvalidParameters = "INVALID_PARAMETERS"

	// Search
	CodeEmptySearchQuery = "EMPTY_SEARCH_QUERY"
)

// ErrorResponse represents the standard error response structure
type ErrorResponse struct {
	Code    string      `json:"code"`
	Message string      `json:"message"`
	Details interface{} `json:"details,omitempty"`
}

// ValidationErrorMessage represents a structured validation error message
type ValidationErrorMessage struct {
	Field      string `json:"field"`
	Tag        string `json:"tag"`
	Value      string `json:"value,omitempty"`
	Message    string `json:"message"`
	Suggestion string `json:"suggestion,omitempty"`
}
