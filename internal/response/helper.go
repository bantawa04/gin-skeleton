package response

import (
	"net/http"

	validators "gin/internal/validator"

	"github.com/gin-gonic/gin"
)

// ResponseHelper provides methods for standardized API responses
type ResponseHelper struct{}

// PaginatedResponse sends a successful paginated response
func (r *ResponseHelper) PaginatedResponse(c *gin.Context, data interface{}, message string) {
	// Check if data is a map with pagination info
	paginatedData, ok := data.(map[string]interface{})
	if !ok {
		// If not, fall back to regular response
		r.OkResponse(c, data, message)
		return
	}

	// Create response with pagination metadata
	response := gin.H{
		"success": true,
		"message": message,
		"data":    paginatedData["data"],
		"meta": gin.H{
			"total":        paginatedData["total"],
			"per_page":     paginatedData["per_page"],
			"current_page": paginatedData["current_page"],
			"last_page":    paginatedData["last_page"],
		},
	}

	c.JSON(http.StatusOK, response)
}

// NewResponseHelper creates a new ResponseHelper
func NewResponseHelper() *ResponseHelper {
	return &ResponseHelper{}
}

// SendPaginatedResponse sends a paginated response
func (h *ResponseHelper) SendPaginatedResponse(c *gin.Context, data interface{}, message string, page, totalPages, perPage, totalItems int) {
	SendPaginatedResponse(c, data, message)
}

// SendResponse sends a success response with data
func (h *ResponseHelper) SendResponse(c *gin.Context, data interface{}, message string, statusCode int) {
	SendResponse(c, data, message, statusCode)
}

// SendError sends an error response
func (h *ResponseHelper) SendError(c *gin.Context, message string, description string, code int, data ...interface{}) {
	SendError(c, message, description, code, data...)
}

// SendSuccess sends a simple success message
func (h *ResponseHelper) SendSuccess(c *gin.Context, message string, statusCode int) {
	SendSuccess(c, message, statusCode)
}

// OkResponse is a shorthand for sending a 200 OK response
func (h *ResponseHelper) OkResponse(c *gin.Context, data interface{}, message string) {
	SendResponse(c, data, message, HTTPOk)
}

// CreatedResponse is a shorthand for sending a 201 Created response
func (h *ResponseHelper) CreatedResponse(c *gin.Context, data interface{}, message string) {
	SendResponse(c, data, message, HTTPCreated)
}

// ValidationError sends a validation error response
func (r *ResponseHelper) ValidationError(c *gin.Context, errors []validators.ValidationError, message string) {
	response := gin.H{
		"success": false,
		"message": message,
		"errors":  errors,
	}
	c.JSON(http.StatusUnprocessableEntity, response)
}
