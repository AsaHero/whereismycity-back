package outerr

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/AsaHero/whereismycity/internal/inerr"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

// HandleError is the main error handling function that integrates with internal_errors
func HandleError(c *gin.Context, err error) {
	var validationErrors validator.ValidationErrors

	switch {
	case errors.As(err, &validationErrors):
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Code:    CodeValidation,
			Message: "Validation failed",
			Details: formatValidationErrors(validationErrors),
		})
		return
	case inerr.IsErrNotFound(err):
		c.JSON(http.StatusNotFound, ErrorResponse{
			Code:    CodeNotFound,
			Message: err.Error(),
		})

	case inerr.IsErrConflict(err):
		c.JSON(http.StatusConflict, ErrorResponse{
			Code:    CodeConflict,
			Message: err.Error(),
		})

	case inerr.IsErrNoChanges(err):
		c.JSON(http.StatusNotModified, ErrorResponse{
			Code:    CodeNoChanges,
			Message: err.Error(),
		})
	default:
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Code:    CodeInternalError,
			Message: err.Error(),
		})
	}
}

// Helper functions for direct error responses
func BadRequest(c *gin.Context, message string) {
	c.JSON(http.StatusBadRequest, ErrorResponse{
		Code:    CodeBadRequest,
		Message: message,
	})
}

// Helper functions for direct error responses
func Internal(c *gin.Context, message string) {
	c.JSON(http.StatusBadRequest, ErrorResponse{
		Code:    CodeInternalError,
		Message: message,
	})
}

func Unauthorized(c *gin.Context, message string) {
	c.JSON(http.StatusUnauthorized, ErrorResponse{
		Code:    CodeUnauthorized,
		Message: message,
	})
}

func Forbidden(c *gin.Context, message string) {
	c.JSON(http.StatusForbidden, ErrorResponse{
		Code:    CodeForbidden,
		Message: message,
	})
}

func TooManyRequests(c *gin.Context, message string) {
	c.JSON(http.StatusUnauthorized, ErrorResponse{
		Code:    CodeTooManyRequests,
		Message: message,
	})
}

func NotFound(c *gin.Context, message string) {
	c.JSON(http.StatusNotFound, ErrorResponse{
		Code:    CodeNotFound,
		Message: message,
	})
}

// getValidationErrorMessage returns detailed validation error information
func getValidationErrorMessage(err validator.FieldError) ValidationErrorMessage {
	var msg ValidationErrorMessage
	msg.Field = err.Field()
	msg.Tag = err.Tag()
	msg.Value = err.Param()

	switch err.Tag() {
	case "required":
		msg.Message = fmt.Sprintf("The %s field is required", err.Field())
		msg.Suggestion = "Please provide a value for this field"

	case "email":
		msg.Message = "Invalid email address format"
		msg.Suggestion = "Please provide a valid email address (e.g., user@example.com)"

	case "min":
		msg.Message = fmt.Sprintf("The %s field must be at least %s characters long", err.Field(), err.Param())
		msg.Suggestion = fmt.Sprintf("Please provide a value with at least %s characters", err.Param())

	case "max":
		msg.Message = fmt.Sprintf("The %s field must not exceed %s characters", err.Field(), err.Param())
		msg.Suggestion = fmt.Sprintf("Please provide a value with no more than %s characters", err.Param())

	case "len":
		msg.Message = fmt.Sprintf("The %s field must be exactly %s characters long", err.Field(), err.Param())
		msg.Suggestion = fmt.Sprintf("Please provide a value that is exactly %s characters long", err.Param())

	case "numeric":
		msg.Message = fmt.Sprintf("The %s field must contain only numbers", err.Field())
		msg.Suggestion = "Please provide a numeric value"

	case "alpha":
		msg.Message = fmt.Sprintf("The %s field must contain only letters", err.Field())
		msg.Suggestion = "Please provide alphabetic characters only"

	case "alphanum":
		msg.Message = fmt.Sprintf("The %s field must contain only letters and numbers", err.Field())
		msg.Suggestion = "Please provide alphanumeric characters only"

	case "url":
		msg.Message = "Invalid URL format"
		msg.Suggestion = "Please provide a valid URL (e.g., https://example.com)"

	case "datetime":
		msg.Message = "Invalid datetime format"
		msg.Suggestion = "Please provide a valid datetime (e.g., 2006-01-02T15:04:05Z)"

	case "uuid":
		msg.Message = "Invalid UUID format"
		msg.Suggestion = "Please provide a valid UUID"

	case "unique":
		msg.Message = fmt.Sprintf("The %s field must be unique", err.Field())
		msg.Suggestion = "Please provide a unique value"

	case "oneof":
		msg.Message = fmt.Sprintf("The %s field must be one of: %s", err.Field(), err.Param())
		msg.Suggestion = fmt.Sprintf("Please choose one of the allowed values: %s", err.Param())

	default:
		msg.Message = fmt.Sprintf("The %s field failed validation", err.Field())
		msg.Suggestion = "Please check the value and try again"
	}

	return msg
}

// formatValidationErrors formats validator.ValidationErrors into a structured format
func formatValidationErrors(errs validator.ValidationErrors) []ValidationErrorMessage {
	validationErrors := make([]ValidationErrorMessage, 0, len(errs))
	for _, err := range errs {
		validationErrors = append(validationErrors, getValidationErrorMessage(err))
	}
	return validationErrors
}
