package response

import (
	"net/http"
	"strings"

	validators "gin/internal/validator"

	"github.com/gin-gonic/gin"
	"github.com/iancoleman/strcase"
)

// HTTP status codes
const (
	HTTPOk                  = http.StatusOK
	HTTPCreated             = http.StatusCreated
	HTTPUnprocessableEntity = http.StatusUnprocessableEntity
)

// Response represents a standard API response
type Response struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Message string      `json:"message"`
}

// ErrorResponse represents an error API response (Laravel-style)
type ErrorResponse struct {
	Message string              `json:"message"`
	Errors  map[string][]string `json:"errors,omitempty"` // Laravel-style: field -> array of error messages
}

// SendResponse sends a success response with data
func SendResponse(c *gin.Context, data interface{}, message string, statusCode ...int) {
	response := Response{
		Success: true,
		Data:    data,
		Message: message,
	}

	code := HTTPOk // Default to 200 OK
	if len(statusCode) > 0 {
		code = statusCode[0]
	}

	c.JSON(code, response)
}

// SendError sends an error response (Laravel-style)
func SendError(c *gin.Context, message string, description string, code int, data ...interface{}) {
	response := ErrorResponse{
		Message: message,
	}

	// If validation errors are provided, format them in Laravel style
	if len(data) > 0 {
		if validationErrors, ok := data[0].([]validators.ValidationError); ok {
			errors := make(map[string][]string)
			for _, ve := range validationErrors {
				// Convert field name to snake_case for consistency
				field := toSnakeCase(ve.Field)
				errors[field] = append(errors[field], ve.Message)
			}
			response.Errors = errors
		}
	}

	c.JSON(code, response)
}

// toSnakeCase converts a string to snake_case
func toSnakeCase(s string) string {
	// Use strcase for proper conversion
	return strings.ToLower(strcase.ToSnake(s))
}

// SendSuccess sends a simple success message
func SendSuccess(c *gin.Context, message string, statusCode int) {
	response := Response{
		Success: true,
		Message: message,
	}

	c.JSON(statusCode, response)
}

// ValidationError sends a validation error response (Laravel-style)
func ValidationError(c *gin.Context, errors []validators.ValidationError, message string) {
	errorMap := make(map[string][]string)
	for _, ve := range errors {
		field := toSnakeCase(ve.Field)
		errorMap[field] = append(errorMap[field], ve.Message)
	}

	response := ErrorResponse{
		Message: message,
		Errors:  errorMap,
	}
	c.JSON(http.StatusUnprocessableEntity, response)
}
