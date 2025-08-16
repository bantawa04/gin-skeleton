package user

import (
	"context"
	"errors"
	exceptions "gin/internal/api/exception"
	"gin/internal/models"
	repository "gin/internal/repository/user"
	"gin/internal/request"
	"gin/internal/utils"

	"golang.org/x/crypto/bcrypt"
)

// UserService implements UserServiceInterface
type UserService struct {
	userRepo *repository.UserRepository
}

// NewUserService creates a new user service
func NewUserService(userRepo *repository.UserRepository) UserServiceInterface {
	return &UserService{
		userRepo: userRepo,
	}
}

// GetAllUsers retrieves all users
func (s *UserService) GetAllUsers(ctx context.Context) ([]*models.User, error) {
	return s.userRepo.GetAll(ctx)
}

// GetAllUsersPaginated retrieves users with pagination
func (s *UserService) GetAllUsersPaginated(ctx context.Context, page, perPage int) ([]*models.User, int64, error) {
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
func (s *UserService) GetUserByID(ctx context.Context, id string) (*models.User, error) {
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

// CreateUser creates a new user
func (s *UserService) CreateUser(ctx context.Context, request request.UserCreateRequest) (interface{}, error) {
	// Check if user already exists
	existingUser, err := s.userRepo.FindByEmail(ctx, request.Email)
	if err != nil {
		return nil, err
	}

	if existingUser != nil {
		return nil, exceptions.ValidationError("User already exists", nil, nil)
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &models.User{
		FirstName: &request.FirstName,
		LastName:  &request.LastName,
		Email:     request.Email,
		Password:  string(hashedPassword),
	}

	return s.userRepo.Create(ctx, user)
}

// UpdateUser updates an existing user
func (s *UserService) UpdateUser(ctx context.Context, request request.UserUpdateRequest, id string) (interface{}, error) {
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

	updates, err := utils.MapStructToUpdates(request)
	if err != nil {
		return nil, err
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

func (s *UserService) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	return s.userRepo.FindByEmail(ctx, email)
}
