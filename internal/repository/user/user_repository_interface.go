package repository

import (
	"context"
	"gin/internal/models"
	"gin/internal/request"
)

// UserRepositoryInterface defines user repository operations
type UserRepositoryInterface interface {
	// Basic CRUD operations
	GetAll(ctx context.Context) ([]*models.User, error)
	Create(ctx context.Context, request request.UserCreateRequest) error
	Update(ctx context.Context, user *models.User) error
	Delete(ctx context.Context, id string) error

	// Find operations
	FindByID(ctx context.Context, id string) (*models.User, error)
	FindByEmail(ctx context.Context, email string) (*models.User, error)
}
