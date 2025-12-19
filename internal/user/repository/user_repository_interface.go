package repository

import (
	"context"
	"gin/internal/user"
)

// UserRepositoryInterface defines user repository operations
type UserRepositoryInterface interface {
	// Basic CRUD operations
	GetAll(ctx context.Context) ([]*user.User, error)
	GetAllPaginated(ctx context.Context, page, perPage int) ([]*user.User, int64, error)
	Create(ctx context.Context, user *user.User) (*user.User, error)
	Update(ctx context.Context, user *user.User) error
	UpdateFields(ctx context.Context, id string, updates map[string]interface{}) error
	Delete(ctx context.Context, id string) error

	// Find operations
	FindByID(ctx context.Context, id string) (*user.User, error)
	FindByEmail(ctx context.Context, email string) (*user.User, error)
}
