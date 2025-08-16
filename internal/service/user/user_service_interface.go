package user

import (
	"context"
	"gin/internal/models"
	"gin/internal/request"
)

type UserServiceInterface interface {
	GetAllUsers(ctx context.Context) ([]*models.User, error)
	GetAllUsersPaginated(ctx context.Context, page, perPage int) ([]*models.User, int64, error)
	GetUserByID(ctx context.Context, id string) (*models.User, error)
	CreateUser(ctx context.Context, request request.UserCreateRequest) (interface{}, error)
	UpdateUser(ctx context.Context, request request.UserUpdateRequest, id string) (interface{}, error)
	DeleteUser(ctx context.Context, id string) error
	GetUserByEmail(ctx context.Context, email string) (*models.User, error)
}
