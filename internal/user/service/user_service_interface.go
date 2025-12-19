package service

import (
	"context"
	"gin/internal/user"
)

type UserServiceInterface interface {
	GetAllUsers(ctx context.Context) ([]*user.User, error)
	GetAllUsersPaginated(ctx context.Context, page, perPage int) ([]*user.User, int64, error)
	GetUserByID(ctx context.Context, id string) (*user.User, error)
	CreateUser(ctx context.Context, req user.SignupInput) (*user.User, error)
	UpdateUser(ctx context.Context, updates map[string]interface{}, password *string, id string) (*user.User, error)
	DeleteUser(ctx context.Context, id string) error
	GetUserByEmail(ctx context.Context, email string) (*user.User, error)
}
