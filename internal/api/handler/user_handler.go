package handler

import (
	"net/http"
	"strconv"

	exceptions "gin/internal/api/exception"
	"gin/internal/models"
	request "gin/internal/request"
	response "gin/internal/response"
	responseDTO "gin/internal/response"
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
	// Get pagination parameters from query string
	page := 1
	perPage := 10

	if pageStr := c.Query("page"); pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			page = p
		}
	}

	if perPageStr := c.Query("per_page"); perPageStr != "" {
		if pp, err := strconv.Atoi(perPageStr); err == nil && pp > 0 && pp <= 100 {
			perPage = pp
		}
	}

	users, total, err := h.userService.GetAllUsersPaginated(c.Request.Context(), page, perPage)
	if err != nil {
		_ = c.Error(err)
		return
	}

	userResponses := responseDTO.ToPaginatedUserResponse(users, page, perPage, total)
	response.SendResponse(c, userResponses, "users retrieved successfully")
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
		_ = c.Error(err)
		return
	}

	userResponse := responseDTO.ToUserResponse(user)
	response.SendResponse(c, userResponse, "user retrieved successfully")
}

// CreateUser handles POST /users request
func (h *UserHandler) CreateUser(c *gin.Context) {
	var request request.UserCreateRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		errMsg := err.Error()
		appErr := exceptions.ValidationError("Invalid request format", &errMsg)
		_ = c.Error(appErr)
		return
	}

	if err := h.validator.Struct(request); err != nil {
		validationErrors := h.validator.GenerateValidationErrors(err)

		appErr := exceptions.ValidationError("Validation failed", nil, validationErrors)
		_ = c.Error(appErr)
		return
	}

	createdUser, err := h.userService.CreateUser(c.Request.Context(), request)
	if err != nil {
		_ = c.Error(err)
		return
	}

	userResponse := responseDTO.ToUserResponse(createdUser.(*models.User))
	response.SendResponse(c, userResponse, "user created successfully")
}

// UpdateUser handles PUT /users/:id request
func (h *UserHandler) UpdateUser(c *gin.Context) {
	id := c.Param("id")
	var request request.UserUpdateRequest

	if id == "" {
		appErr := exceptions.ValidationError("User ID is required", nil, nil)
		_ = c.Error(appErr)
		return
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		errMsg := err.Error()
		appErr := exceptions.ValidationError("Invalid request format", &errMsg)
		_ = c.Error(appErr)
		return
	}

	if err := h.validator.Struct(request); err != nil {
		validationErrors := h.validator.GenerateValidationErrors(err)

		appErr := exceptions.ValidationError("Validation failed", nil, validationErrors)
		_ = c.Error(appErr)
		return
	}

	updatedUser, err := h.userService.UpdateUser(c.Request.Context(), request, id)
	if err != nil {
		_ = c.Error(err)
		return
	}

	userResponse := responseDTO.ToUserResponse(updatedUser.(*models.User))
	response.SendResponse(c, userResponse, "user updated successfully")
}

// DeleteUser handles DELETE /users/:id request
func (h *UserHandler) DeleteUser(c *gin.Context) {
	id := c.Param("id")

	if id == "" {
		appErr := exceptions.ValidationError("User ID is required", nil, nil)
		_ = c.Error(appErr)
		return
	}
	err := h.userService.DeleteUser(c.Request.Context(), id)
	if err != nil {
		_ = c.Error(err)
		return
	}

	response.SendResponse(c, nil, "user deleted successfully")
}
