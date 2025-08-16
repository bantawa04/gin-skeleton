package user

import (
	"context"
	"gin/internal/models"
	"gin/internal/request"
)

// UserServiceInterface defines user business logic operations
type UserServiceInterface interface {
	// User management
	GetAllUsers(ctx context.Context) ([]*models.User, error)
	GetUserByID(ctx context.Context, id string) (*models.User, error)
	CreateUser(ctx context.Context, request request.UserCreateRequest) (interface{}, error)
	UpdateUser(ctx context.Context, user *models.User) (interface{}, error)
	DeleteUser(ctx context.Context, id string) error
}
