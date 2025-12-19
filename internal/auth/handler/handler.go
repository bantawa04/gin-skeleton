package handler

import (
	"net/http"
	"time"

	"gin/internal/auth"
	"gin/internal/constant"
	exceptions "gin/internal/exception"
	refreshtoken "gin/internal/refresh_token"
	refreshsvc "gin/internal/refresh_token/service"
	"gin/internal/response"
	userdomain "gin/internal/user"
	usersvc "gin/internal/user/service"
	"gin/internal/utils"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type AuthHandler struct {
	userService         usersvc.UserServiceInterface
	jwtManager          *utils.JWTManager
	refreshTokenService refreshsvc.RefreshTokenServiceInterface
}

func NewAuthHandler(userService usersvc.UserServiceInterface, jwtManager *utils.JWTManager, refreshTokenService refreshsvc.RefreshTokenServiceInterface) *AuthHandler {
	return &AuthHandler{
		userService:         userService,
		jwtManager:          jwtManager,
		refreshTokenService: refreshTokenService,
	}
}

// Signup creates a new user account with email and name only
// User will receive a verification email to set their password
func (h *AuthHandler) Signup(c *gin.Context) {
	var req auth.SignupRequest

	// Bind and validate JSON request
	if err := c.ShouldBindJSON(&req); err != nil {
		// Try to extract validation errors from binding error
		validationErrors := utils.ExtractBindingErrors(err)
		if len(validationErrors) > 0 {
			appErr := exceptions.ValidationError("The given data was invalid.", nil, validationErrors)
			_ = c.Error(appErr)
			return
		}
		// If no validation errors extracted, it's a JSON parsing error
		errMsg := "Invalid request format. Please check your JSON syntax."
		appErr := exceptions.ValidationError(errMsg, nil)
		_ = c.Error(appErr)
		return
	}

	// Create user (without password - password will be set after email verification)
	signupInput := userdomain.SignupInput{
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Email:     req.Email,
	}
	_, err := h.userService.CreateUser(c.Request.Context(), signupInput)
	if err != nil {
		_ = c.Error(err)
		return
	}

	// TODO: Send verification email here
	// The email should contain a link to the password setup page
	response.SendSuccess(c, "User created successfully. Please check your email to verify your account and set your password.", http.StatusCreated)
}

// Login authenticates a user and returns JWT tokens
// @Summary      User login
// @Description  Authenticate user with email and password, receive access and refresh tokens
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        credentials  body      auth.LoginRequest  true  "Login credentials"
// @Success      200          {object}  response.Response{data=auth.LoginResponseDTO}
// @Failure      400          {object}  response.ErrorResponse
// @Failure      401          {object}  response.ErrorResponse
// @Failure      404          {object}  response.ErrorResponse
// @Failure      422          {object}  response.ErrorResponse
// @Failure      500          {object}  response.ErrorResponse
// @Router       /auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var req auth.LoginRequest

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

	u, err := h.userService.GetUserByEmail(c.Request.Context(), req.Email)

	if err != nil {
		_ = c.Error(err)
		return
	}

	if u == nil {
		appErr := exceptions.NotFoundError("User not found with given email", nil, nil)
		_ = c.Error(appErr)
		return
	}

	// Check if user is active
	if u.Status != constant.UserStatusActive {
		appErr := exceptions.UnauthorizedError("Your account is not active. Please verify your email and set your password", nil, nil)
		_ = c.Error(appErr)
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(req.Password)); err != nil {
		appErr := exceptions.UnauthorizedError("Invalid credentials", nil, nil)
		_ = c.Error(appErr)
		return
	}

	// Generate access token
	accessToken, err := h.jwtManager.GenerateAccessToken(u.ID)
	if err != nil {
		appErr := exceptions.InternalError("Failed to generate access token", nil, nil)
		_ = c.Error(appErr)
		return
	}

	// Generate refresh token
	refreshToken, err := h.jwtManager.GenerateRefreshToken(u.ID)
	if err != nil {
		appErr := exceptions.InternalError("Failed to generate refresh token", nil, nil)
		_ = c.Error(appErr)
		return
	}

	// Save refresh token to database
	refreshTokenModel := &refreshtoken.RefreshToken{
		UserID:    u.ID,
		Token:     refreshToken,
		ExpiresAt: time.Now().Add(h.jwtManager.GetRefreshExpiry()),
		Revoked:   false,
	}

	_, err = h.refreshTokenService.Create(c.Request.Context(), refreshTokenModel)
	if err != nil {
		appErr := exceptions.InternalError("Failed to save refresh token", nil, nil)
		_ = c.Error(appErr)
		return
	}

	// Create response DTO
	loginResponseDTO := auth.LoginResponseDTO{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		TokenType:    "Bearer",
		ExpiresIn:    int64(h.jwtManager.GetAccessExpiry().Seconds()),
		User:         userdomain.FromUserModel(*u),
	}

	// Send success response
	response.SendResponse(c, loginResponseDTO, "Login successful")
}

// RefreshToken generates a new access token using a valid refresh token
// @Summary      Refresh access token
// @Description  Generate a new access token using a valid refresh token
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        refresh  body      auth.RefreshTokenRequest  true  "Refresh token"
// @Success      200      {object}  response.Response{data=auth.RefreshTokenResponseDTO}
// @Failure      400      {object}  response.ErrorResponse
// @Failure      401      {object}  response.ErrorResponse
// @Failure      422      {object}  response.ErrorResponse
// @Router       /auth/refresh [post]
func (h *AuthHandler) RefreshToken(c *gin.Context) {
	var req auth.RefreshTokenRequest

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

	// Validate refresh token exists in database and is not revoked
	dbRefreshToken, err := h.refreshTokenService.FindByToken(c.Request.Context(), req.RefreshToken)
	if err != nil {
		appErr := exceptions.InternalError("Failed to validate refresh token", nil, nil)
		_ = c.Error(appErr)
		return
	}

	if dbRefreshToken == nil || !dbRefreshToken.IsValid() {
		appErr := exceptions.UnauthorizedError("Invalid or revoked refresh token", nil, nil)
		_ = c.Error(appErr)
		return
	}

	// Validate JWT signature and extract user ID
	claims, err := h.jwtManager.ValidateToken(req.RefreshToken)
	if err != nil {
		appErr := exceptions.UnauthorizedError("Invalid refresh token", nil, nil)
		_ = c.Error(appErr)
		return
	}

	// Ensure this is a refresh token
	if claims.Type != "refresh" {
		appErr := exceptions.UnauthorizedError("Invalid token type, refresh token required", nil, nil)
		_ = c.Error(appErr)
		return
	}

	// ROTATION STEP 1: Revoke the old refresh token
	err = h.refreshTokenService.RevokeByToken(c.Request.Context(), req.RefreshToken)
	if err != nil {
		appErr := exceptions.InternalError("Failed to revoke old refresh token", nil, nil)
		_ = c.Error(appErr)
		return
	}

	// ROTATION STEP 2: Generate new access token
	accessToken, err := h.jwtManager.GenerateAccessToken(claims.UserID)
	if err != nil {
		appErr := exceptions.InternalError("Failed to generate access token", nil, nil)
		_ = c.Error(appErr)
		return
	}

	// ROTATION STEP 3: Generate new refresh token
	newRefreshToken, err := h.jwtManager.GenerateRefreshToken(claims.UserID)
	if err != nil {
		appErr := exceptions.InternalError("Failed to generate refresh token", nil, nil)
		_ = c.Error(appErr)
		return
	}

	// ROTATION STEP 4: Save new refresh token to database
	newRefreshTokenModel := &refreshtoken.RefreshToken{
		UserID:    claims.UserID,
		Token:     newRefreshToken,
		ExpiresAt: time.Now().Add(h.jwtManager.GetRefreshExpiry()),
		Revoked:   false,
	}

	_, err = h.refreshTokenService.Create(c.Request.Context(), newRefreshTokenModel)
	if err != nil {
		appErr := exceptions.InternalError("Failed to save new refresh token", nil, nil)
		_ = c.Error(appErr)
		return
	}

	// Create response DTO with both new tokens
	refreshResponseDTO := auth.RefreshTokenResponseDTO{
		AccessToken:  accessToken,
		RefreshToken: newRefreshToken, // Return new refresh token
		TokenType:    "Bearer",
		ExpiresIn:    int64(h.jwtManager.GetAccessExpiry().Seconds()),
	}

	// Send success response
	response.SendResponse(c, refreshResponseDTO, "Token refreshed successfully")
}

// Logout revokes all refresh tokens for the authenticated user
// @Summary      User logout
// @Description  Revoke all refresh tokens for the authenticated user, effectively logging them out. Requires JWT authentication via Authorization header.
// @Tags         auth
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Success      200          {object}  response.Response
// @Failure      401          {object}  response.ErrorResponse
// @Failure      500          {object}  response.ErrorResponse
// @Router       /auth/logout [post]
func (h *AuthHandler) Logout(c *gin.Context) {
	// Get user ID from context (set by JWTAuthMiddleware)
	userID, err := utils.RequireUserID(c)
	if err != nil {
		appErr := exceptions.UnauthorizedError("User ID not found in context", nil, nil)
		_ = c.Error(appErr)
		return
	}

	// Revoke all refresh tokens for this user
	err = h.refreshTokenService.RevokeAllUserTokens(c.Request.Context(), userID)
	if err != nil {
		appErr := exceptions.InternalError("Failed to revoke refresh tokens", nil, nil)
		_ = c.Error(appErr)
		return
	}

	response.SendResponse(c, nil, "Logout successful")
}
