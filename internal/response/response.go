package response

import (
	"net/http"

	"gin/internal/validator"
	"github.com/gin-gonic/gin"
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

// ErrorResponse represents an error API response
type ErrorResponse struct {
	Success     bool        `json:"success"`
	Message     string      `json:"message"`
	Description string      `json:"description,omitempty"`
	Data        interface{} `json:"data,omitempty"`
}

// SendPaginatedResponse sends a paginated response
func SendPaginatedResponse(c *gin.Context, data interface{}, message string) {

	paginatedData, ok := data.(map[string]interface{})
	if !ok {
		// If not, fall back to regular response
		SendResponse(c, data, message)
		return
	}

	// Create response with pagination metadata
	response := gin.H{
		"success": true,
		"message": message,
		"data":    paginatedData["data"],
		"meta": gin.H{
			"page":       paginatedData["page"],
			"totalPages": paginatedData["totalPages"],
			"perPage":    paginatedData["perPage"],
			"totalItems": paginatedData["totalItems"],
		},
	}

	c.JSON(http.StatusOK, response)
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

// SendError sends an error response
func SendError(c *gin.Context, message string, description string, code int, data ...interface{}) {
	response := ErrorResponse{
		Success:     false,
		Message:     message,
		Description: description,
	}

	if len(data) > 0 {
		response.Data = data[0]
	}

	c.JSON(code, response)
}

// SendSuccess sends a simple success message
func SendSuccess(c *gin.Context, message string, statusCode int) {
	response := Response{
		Success: true,
		Message: message,
	}

	c.JSON(statusCode, response)
}

func ValidationError(c *gin.Context, errors []validators.ValidationError, message string) {
	response := gin.H{
		"success": false,
		"message": message,
		"errors":  errors,
	}
	c.JSON(http.StatusUnprocessableEntity, response)
}
