package handler

import (
	"net/http"
	"strconv"

	exceptions "gin/internal/exception"
	"gin/internal/response"
	"gin/internal/user"
	usersvc "gin/internal/user/service"
	"gin/internal/utils"

	"github.com/gin-gonic/gin"
)

// GetAllUsers handles GET /users request
// @Summary      List all users
// @Description  Get a paginated list of all users
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        page      query     int  false  "Page number"  default(1)
// @Param        per_page  query     int  false  "Items per page"  default(10)  maximum(100)
// @Success      200       {object}  response.Response{data=user.PaginatedUserDTO}
// @Failure      500       {object}  response.ErrorResponse
// @Router       /users [get]

// UserHandler handles HTTP requests for user operations
type UserHandler struct {
	userService usersvc.UserServiceInterface
}

// NewUserHandler creates a new user handler
func NewUserHandler(userService usersvc.UserServiceInterface) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

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

	paginatedDTO := user.ToPaginatedUserDTO(users, page, perPage, total)
	response.SendResponse(c, paginatedDTO, "users retrieved successfully")
}

// GetUserByID handles GET /users/:id request
// @Summary      Get user by ID
// @Description  Get a specific user by their ID
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "User ID"
// @Success      200  {object}  response.Response{data=user.UserDTO}
// @Failure      400  {object}  response.ErrorResponse
// @Failure      404  {object}  response.ErrorResponse
// @Failure      500  {object}  response.ErrorResponse
// @Router       /users/{id} [get]
func (h *UserHandler) GetUserByID(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user ID is required"})
		return
	}

	u, err := h.userService.GetUserByID(c.Request.Context(), id)
	if err != nil {
		_ = c.Error(err)
		return
	}

	userDTO := user.FromUserModel(*u)
	response.SendResponse(c, userDTO, "user retrieved successfully")
}

// UpdateUser handles PUT /users/:id request
// @Summary      Update user
// @Description  Update an existing user's information
// @Tags         users
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id    path      string                   true  "User ID"
// @Param        user  body      user.UserUpdateRequest  true  "User update data"
// @Success      200   {object}  response.Response{data=user.UserDTO}
// @Failure      400   {object}  response.ErrorResponse
// @Failure      401   {object}  response.ErrorResponse
// @Failure      404   {object}  response.ErrorResponse
// @Failure      422   {object}  response.ErrorResponse
// @Failure      500   {object}  response.ErrorResponse
// @Router       /users/{id} [put]
func (h *UserHandler) UpdateUser(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		appErr := exceptions.ValidationError("User ID is required", nil, nil)
		_ = c.Error(appErr)
		return
	}

	var req user.UserUpdateRequest

	// Bind and validate JSON request
	if err := c.ShouldBindJSON(&req); err != nil {
		validationErrors := utils.ExtractBindingErrors(err)
		if len(validationErrors) > 0 {
			appErr := exceptions.ValidationError("The given data was invalid.", nil, validationErrors)
			_ = c.Error(appErr)
			return
		}
		errMsg := "Invalid request format. Please check your JSON syntax."
		appErr := exceptions.ValidationError(errMsg, nil)
		_ = c.Error(appErr)
		return
	}

	// Convert request to update map
	updates := make(map[string]interface{})
	if req.Name != nil {
		updates["first_name"] = *req.Name
	}
	if req.Email != nil {
		updates["email"] = *req.Email
	}

	updatedUser, err := h.userService.UpdateUser(c.Request.Context(), updates, req.Password, id)
	if err != nil {
		_ = c.Error(err)
		return
	}

	// Convert model to response DTO
	userDTO := user.FromUserModel(*updatedUser)
	response.SendResponse(c, userDTO, "user updated successfully")
}

// DeleteUser handles DELETE /users/:id request
// @Summary      Delete user
// @Description  Delete a user account
// @Tags         users
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      string  true  "User ID"
// @Success      200  {object}  response.Response
// @Failure      400  {object}  response.ErrorResponse
// @Failure      401  {object}  response.ErrorResponse
// @Failure      404  {object}  response.ErrorResponse
// @Failure      500  {object}  response.ErrorResponse
// @Router       /users/{id} [delete]
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
