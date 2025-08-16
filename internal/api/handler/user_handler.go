package handler

import (
	"net/http"

	middleware "gin/internal/api/middleware"
	"gin/internal/models"
	request "gin/internal/request"
	response "gin/internal/response"
	"gin/internal/service/user"
	validators "gin/internal/validator"

	"github.com/gin-gonic/gin"
)

// UserHandler handles HTTP requests for user operations
type UserHandler struct {
	userService user.UserServiceInterface
	validator   *validators.Validator
}

// NewUserHandler creates a new user handler
func NewUserHandler(userService user.UserServiceInterface, validator *validators.Validator) *UserHandler {
	return &UserHandler{
		userService: userService,
		validator:   validator,
	}
}

// GetAllUsers handles GET /users request
func (h *UserHandler) GetAllUsers(c *gin.Context) {
	users, err := h.userService.GetAllUsers(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	response.SendResponse(c, users, "users retrieved successfully")
}

// GetUserByID handles GET /users/:id request
func (h *UserHandler) GetUserByID(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user ID is required"})
		return
	}

	user, err := h.userService.GetUserByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if user == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}

	response.SendResponse(c, user, "user retrieved successfully")
}

// CreateUser handles POST /users request
func (h *UserHandler) CreateUser(c *gin.Context) {
	var request request.UserCreateRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		appErr := middleware.NewValidationError("Invalid request format", err.Error())
		_ = c.Error(appErr)
		return
	}

	// Validate the request
	if err := h.validator.Struct(request); err != nil {
		validationErrors := h.validator.GenerateValidationErrors(err)
		appErr := middleware.NewValidationError("Validation failed", "Request validation failed", validationErrors)
		_ = c.Error(appErr)
		return
	}

	createdUser, err := h.userService.CreateUser(c.Request.Context(), request)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response.SendResponse(c, createdUser, "user created successfully")
}

// UpdateUser handles PUT /users/:id request
func (h *UserHandler) UpdateUser(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user ID is required"})
		return
	}

	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	// Set the ID from URL parameter
	user.ID = id

	updatedUser, err := h.userService.UpdateUser(c.Request.Context(), &user)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response.SendResponse(c, updatedUser, "user updated successfully")
}

// DeleteUser handles DELETE /users/:id request
func (h *UserHandler) DeleteUser(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user ID is required"})
		return
	}

	err := h.userService.DeleteUser(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response.SendResponse(c, nil, "user deleted successfully")
}
