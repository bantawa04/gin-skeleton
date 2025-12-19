package service

import (
	"context"
	"errors"
	"gin/internal/constant"
	exceptions "gin/internal/exception"
	"gin/internal/user"
	userRepository "gin/internal/user/repository"

	"gin/internal/utils"

	"fmt"

	"golang.org/x/crypto/bcrypt"
)

// UserService implements UserServiceInterface
type UserService struct {
	userRepo *userRepository.UserRepository
}

// NewUserService creates a new user service
func NewUserService(userRepo *userRepository.UserRepository) UserServiceInterface {
	return &UserService{
		userRepo: userRepo,
	}
}

// GetAllUsers retrieves all users
func (s *UserService) GetAllUsers(ctx context.Context) ([]*user.User, error) {
	return s.userRepo.GetAll(ctx)
}

// GetAllUsersPaginated retrieves users with pagination
func (s *UserService) GetAllUsersPaginated(ctx context.Context, page, perPage int) ([]*user.User, int64, error) {
	// Validate pagination parameters
	if page < 1 {
		page = 1
	}
	if perPage < 1 {
		perPage = 10 // Default page size
	}
	if perPage > 100 {
		perPage = 100 // Maximum page size
	}

	return s.userRepo.GetAllPaginated(ctx, page, perPage)
}

// GetUserByID retrieves a user by ID
func (s *UserService) GetUserByID(ctx context.Context, id string) (*user.User, error) {
	if id == "" {
		return nil, errors.New("user ID is required")
	}

	user, err := s.userRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, exceptions.NotFoundError("User not found", nil, nil)
	}

	return user, nil
}

// CreateUser creates a new user during signup
// A random password is generated since password field is required in the database
// User will set their own password after email verification
func (s *UserService) CreateUser(ctx context.Context, req user.SignupInput) (*user.User, error) {
	// Check if user already exists
	existingUser, err := s.userRepo.FindByEmail(ctx, req.Email)
	if err != nil {
		return nil, err
	}

	if existingUser != nil {
		return nil, exceptions.ValidationError("User already exists with this email", nil, nil)
	}

	// Generate a random password since password field is required in the database
	// User will set their own password after email verification
	randomPassword := utils.GeneratePassword()
	fmt.Println("randomPassword", randomPassword)
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(randomPassword), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	// Status is set to 'inactive' by default, indicating user needs to verify email and set password
	user := user.User{
		FirstName: &req.FirstName,
		LastName:  &req.LastName,
		Email:     req.Email,
		Password:  string(hashedPassword),
		Status:    constant.UserStatusInactive, // User is inactive until email is verified and password is set
	}

	return s.userRepo.Create(ctx, &user)
}

// UpdateUser updates an existing user
func (s *UserService) UpdateUser(ctx context.Context, updates map[string]interface{}, password *string, id string) (*user.User, error) {
	if id == "" {
		return nil, errors.New("user ID is required")
	}

	existingUser, err := s.userRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if existingUser == nil {
		return nil, exceptions.NotFoundError("User not found", nil, nil)
	}

	// Handle password separately if provided
	if password != nil && *password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(*password), bcrypt.DefaultCost)
		if err != nil {
			return nil, err
		}
		updates["password"] = string(hashedPassword)
	}

	err = s.userRepo.UpdateFields(ctx, id, updates)
	if err != nil {
		return nil, err
	}

	updatedUser, err := s.userRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return updatedUser, nil
}

// DeleteUser deletes a user by ID
func (s *UserService) DeleteUser(ctx context.Context, id string) error {

	existingUser, err := s.userRepo.FindByID(ctx, id)

	if err != nil {
		return err
	}

	if existingUser == nil {
		return exceptions.NotFoundError("User not found", nil, nil)
	}

	return s.userRepo.Delete(ctx, id)
}

func (s *UserService) GetUserByEmail(ctx context.Context, email string) (*user.User, error) {
	return s.userRepo.FindByEmail(ctx, email)
}
