package user

import (
	"context"
	"errors"
	"gin/internal/models"
	repository "gin/internal/repository/user"
	"gin/internal/request"
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
		return nil, errors.New("user not found")
	}

	return user, nil
}

// GetUserByEmail retrieves a user by email
func (s *UserService) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	if email == "" {
		return nil, errors.New("email is required")
	}

	user, err := s.userRepo.FindByEmail(ctx, email)
	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, errors.New("user not found")
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
		return nil, errors.New("user with this email already exists")
	}

	// Convert request to user model with only the fields provided
	user := &models.User{
		FirstName: &request.FirstName,
		LastName:  &request.LastName,
		Email:     request.Email,
	}

	return s.userRepo.Create(ctx, user)
}

// UpdateUser updates an existing user
func (s *UserService) UpdateUser(ctx context.Context, user *models.User) (interface{}, error) {
	if user == nil {
		return nil, errors.New("user data is required")
	}

	if user.ID == "" {
		return nil, errors.New("user ID is required")
	}

	// Check if user exists
	existingUser, err := s.userRepo.FindByID(ctx, user.ID)
	if err != nil {
		return nil, err
	}

	if existingUser == nil {
		return nil, errors.New("user not found")
	}

	err = s.userRepo.Update(ctx, user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

// DeleteUser deletes a user by ID
func (s *UserService) DeleteUser(ctx context.Context, id string) error {
	if id == "" {
		return errors.New("user ID is required")
	}

	// Check if user exists
	existingUser, err := s.userRepo.FindByID(ctx, id)
	if err != nil {
		return err
	}

	if existingUser == nil {
		return errors.New("user not found")
	}

	return s.userRepo.Delete(ctx, id)
}
