package exception

import (
	"errors"
	response "gin/internal/response"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ErrorType string

const (
	ErrorTypeValidation   ErrorType = "VALIDATION_ERROR"
	ErrorTypeInternal     ErrorType = "INTERNAL_ERROR"
	ErrorTypeNotFound     ErrorType = "NOT_FOUND"
	ErrorTypeUnauthorized ErrorType = "UNAUTHORIZED"
	ErrorTypeForbidden    ErrorType = "FORBIDDEN"
)

type AppError struct {
	Type        ErrorType
	Message     string
	Description *string
	Data        interface{}
}

func (e AppError) Error() string {
	return e.Message
}

func ValidationError(message string, description *string, data ...interface{}) AppError {
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

func InternalError(message string, description *string, data ...interface{}) AppError {
	var errorData interface{}
	if len(data) > 0 {
		errorData = data[0]
	}
	return AppError{
		Type:        ErrorTypeInternal,
		Message:     message,
		Description: description,
		Data:        errorData,
	}
}

func NotFoundError(message string, description *string, data ...interface{}) AppError {
	var errorData interface{}
	if len(data) > 0 {
		errorData = data[0]
	}
	return AppError{
		Type:        ErrorTypeNotFound,
		Message:     message,
		Description: description,
		Data:        errorData,
	}
}

func UnauthorizedError(message string, description *string, data ...interface{}) AppError {
	var errorData interface{}
	if len(data) > 0 {
		errorData = data[0]
	}
	return AppError{
		Type:        ErrorTypeUnauthorized,
		Message:     message,
		Description: description,
		Data:        errorData,
	}
}

func ForbiddenError(message string, description *string, data ...interface{}) AppError {
	var errorData interface{}
	if len(data) > 0 {
		errorData = data[0]
	}
	return AppError{
		Type:        ErrorTypeForbidden,
		Message:     message,
		Description: description,
		Data:        errorData,
	}
}

func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Process request
		c.Next()

		if len(c.Errors) > 0 {
			err := c.Errors.Last().Err

			log.Printf("Error: %v", err)

			var appErr AppError
			if errors.As(err, &appErr) {
				// Handle application errors
				switch appErr.Type {
				case ErrorTypeValidation:
					// Use default Laravel validation message if not provided
					message := appErr.Message
					if message == "" {
						message = "The given data was invalid."
					}
					response.SendError(c, message, "", http.StatusUnprocessableEntity, appErr.Data)
				case ErrorTypeInternal:
					desc := ""
					if appErr.Description != nil {
						desc = *appErr.Description
					}
					response.SendError(c, appErr.Message, desc, http.StatusInternalServerError)
				case ErrorTypeNotFound:
					desc := ""
					if appErr.Description != nil {
						desc = *appErr.Description
					}
					response.SendError(c, appErr.Message, desc, http.StatusNotFound)
				case ErrorTypeUnauthorized:
					desc := ""
					if appErr.Description != nil {
						desc = *appErr.Description
					}
					response.SendError(c, appErr.Message, desc, http.StatusUnauthorized)
				case ErrorTypeForbidden:
					desc := ""
					if appErr.Description != nil {
						desc = *appErr.Description
					}
					response.SendError(c, appErr.Message, desc, http.StatusForbidden)
				default:
					response.SendError(c, "An unexpected error occurred", err.Error(), http.StatusInternalServerError)
				}
			} else {
				// Handle generic errors
				response.SendError(c, "An unexpected error occurred", err.Error(), http.StatusInternalServerError)
			}
		}
	}
}
