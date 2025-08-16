package middleware

import (
	"errors"
	response "gin/internal/response"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

// ErrorType represents different types of errors
type ErrorType string

const (
	// ErrorTypeValidation represents validation errors
	ErrorTypeValidation ErrorType = "VALIDATION_ERROR"
	// ErrorTypeInternal represents internal server errors
	ErrorTypeInternal ErrorType = "INTERNAL_ERROR"
	// ErrorTypeNotFound represents not found errors
	ErrorTypeNotFound ErrorType = "NOT_FOUND"
	// ErrorTypeUnauthorized represents unauthorized errors
	ErrorTypeUnauthorized ErrorType = "UNAUTHORIZED"
	// ErrorTypeForbidden represents forbidden errors
	ErrorTypeForbidden ErrorType = "FORBIDDEN"
)

// AppError represents an application error
type AppError struct {
	Type        ErrorType
	Message     string
	Description string
	Data        interface{}
}

// Error implements the error interface
func (e AppError) Error() string {
	return e.Message
}

// NewValidationError creates a new validation error
func NewValidationError(message string, description string, data ...interface{}) AppError {
	var errorData interface{}
	if len(data) > 0 {
		errorData = data[0]
	}
	return AppError{
		Type:        ErrorTypeValidation,
		Message:     message,
		Description: description,
		Data:        errorData,
	}
}

// NewInternalError creates a new internal error
func NewInternalError(message string, description string) AppError {
	return AppError{
		Type:        ErrorTypeInternal,
		Message:     message,
		Description: description,
	}
}

// NewNotFoundError creates a new not found error
func NewNotFoundError(message string, description string) AppError {
	return AppError{
		Type:        ErrorTypeNotFound,
		Message:     message,
		Description: description,
	}
}

// NewUnauthorizedError creates a new unauthorized error
func NewUnauthorizedError(message string, description string) AppError {
	return AppError{
		Type:        ErrorTypeUnauthorized,
		Message:     message,
		Description: description,
	}
}

// NewForbiddenError creates a new forbidden error
func NewForbiddenError(message string, description string) AppError {
	return AppError{
		Type:        ErrorTypeForbidden,
		Message:     message,
		Description: description,
	}
}

// ErrorHandler is a middleware that handles errors
func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Process request
		c.Next()

		// Check if there are any errors
		if len(c.Errors) > 0 {
			// Get the last error
			err := c.Errors.Last().Err

			// Log the error
			log.Printf("Error: %v", err)

			// Handle different types of errors
			var appErr AppError
			if errors.As(err, &appErr) {
				// Handle application errors
				switch appErr.Type {
				case ErrorTypeValidation:
					response.SendErrorResponse(c, http.StatusUnprocessableEntity, appErr.Message)
				case ErrorTypeInternal:
					response.SendErrorResponse(c, http.StatusInternalServerError, appErr.Message) // Changed from respHelper to responses
				case ErrorTypeNotFound:
					response.SendErrorResponse(c, http.StatusNotFound, appErr.Message)
				case ErrorTypeUnauthorized:
					response.SendErrorResponse(c, http.StatusUnauthorized, appErr.Message)
				case ErrorTypeForbidden:
					response.SendErrorResponse(c, http.StatusForbidden, appErr.Message)
				default:
					response.SendErrorResponse(c, http.StatusInternalServerError, "An unexpected error occurred")
				}
			} else {
				// Handle generic errors
				response.SendErrorResponse(c, http.StatusInternalServerError, "An unexpected error occurred")
			}
		}
	}
}
