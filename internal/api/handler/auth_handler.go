package handler

import (
	exceptions "gin/internal/api/exception"
	request "gin/internal/request"
	response "gin/internal/response"
	"gin/internal/service/user"
	"gin/internal/utils"
	validators "gin/internal/validator"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type AuthHandler struct {
	userService user.UserServiceInterface
	validator   *validators.Validator
	jwtManager  *utils.JWTManager
}

func NewAuthHandler(userService user.UserServiceInterface, validator *validators.Validator, jwtManager *utils.JWTManager) *AuthHandler {
	return &AuthHandler{
		userService: userService,
		validator:   validator,
		jwtManager:  jwtManager,
	}
}

func (h *AuthHandler) Login(c *gin.Context) {
	var request request.LoginRequest

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

	user, err := h.userService.GetUserByEmail(c.Request.Context(), request.Email)

	if err != nil {
		_ = c.Error(err)
		return
	}

	if user == nil {
		appErr := exceptions.NotFoundError("User not found with given email", nil, nil)
		_ = c.Error(appErr)
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(request.Password)); err != nil {
		appErr := exceptions.UnauthorizedError("Invalid credentials", nil, nil)
		_ = c.Error(appErr)
		return
	}

	// Generate access token
	accessToken, err := h.jwtManager.GenerateAccessToken(user.ID)
	if err != nil {
		appErr := exceptions.InternalError("Failed to generate access token", nil, nil)
		_ = c.Error(appErr)
		return
	}

	// Generate refresh token
	refreshToken, err := h.jwtManager.GenerateRefreshToken(user.ID)
	if err != nil {
		appErr := exceptions.InternalError("Failed to generate refresh token", nil, nil)
		_ = c.Error(appErr)
		return
	}

	// Create response
	loginResponse := &response.LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		TokenType:    "Bearer",
		ExpiresIn:    int64(h.jwtManager.GetAccessExpiry().Seconds()),
		User:         response.ToUserResponse(user),
	}

	// Send success response
	response.SendResponse(c, loginResponse, "Login successful")
}

// RefreshToken generates a new access token using a valid refresh token
func (h *AuthHandler) RefreshToken(c *gin.Context) {
	var request request.RefreshTokenRequest

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

	// Generate new access token using refresh token
	accessToken, err := h.jwtManager.RefreshAccessToken(request.RefreshToken)
	if err != nil {
		appErr := exceptions.UnauthorizedError("Invalid refresh token", nil, nil)
		_ = c.Error(appErr)
		return
	}

	// Create response
	refreshResponse := &response.RefreshTokenResponse{
		AccessToken: accessToken,
		TokenType:   "Bearer",
		ExpiresIn:   int64(h.jwtManager.GetAccessExpiry().Seconds()),
	}

	// Send success response
	response.SendResponse(c, refreshResponse, "Token refreshed successfully")
}
