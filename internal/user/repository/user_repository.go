package repository

import (
	"context"
	"gin/internal/user"

	"gorm.io/gorm"
)

// UserRepository handles user database operations
type UserRepository struct {
	db *gorm.DB
}

// NewUserRepository creates a new user repository
func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

// getDB retrieves the database connection from context if transaction exists, otherwise returns default db
func (r *UserRepository) getDB(ctx context.Context) *gorm.DB {
	// Try to get transaction from context (set by transaction middleware)
	if tx, ok := ctx.Value("db_transaction").(*gorm.DB); ok {
		return tx
	}
	return r.db
}

// GetAll retrieves all users
func (r *UserRepository) GetAll(ctx context.Context) ([]*user.User, error) {
	var users []*user.User
	err := r.getDB(ctx).WithContext(ctx).Find(&users).Error
	return users, err
}

// GetAllPaginated retrieves users with pagination
func (r *UserRepository) GetAllPaginated(ctx context.Context, page, perPage int) ([]*user.User, int64, error) {
	var users []*user.User
	var total int64

	// Get total count
	err := r.getDB(ctx).WithContext(ctx).Model(&user.User{}).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	// Calculate offset
	offset := (page - 1) * perPage

	// Get paginated users
	err = r.getDB(ctx).WithContext(ctx).Offset(offset).Limit(perPage).Find(&users).Error
	if err != nil {
		return nil, 0, err
	}

	return users, total, err
}

// Create creates a new user
func (r *UserRepository) Create(ctx context.Context, user *user.User) (*user.User, error) {
	err := r.getDB(ctx).WithContext(ctx).Create(user).Error
	if err != nil {
		return nil, err
	}
	return user, nil
}

// Update updates an existing user
func (r *UserRepository) Update(ctx context.Context, user *user.User) error {
	return r.getDB(ctx).WithContext(ctx).Save(user).Error
}

// UpdateFields updates specific fields of a user
func (r *UserRepository) UpdateFields(ctx context.Context, id string, updates map[string]interface{}) error {
	return r.getDB(ctx).WithContext(ctx).Model(&user.User{}).Where("id = ?", id).Updates(updates).Error
}

// Delete deletes a user by ID
func (r *UserRepository) Delete(ctx context.Context, id string) error {
	return r.getDB(ctx).WithContext(ctx).Where("id = ?", id).Delete(&user.User{}).Error
}

// FindByID finds a user by ID
func (r *UserRepository) FindByID(ctx context.Context, id string) (*user.User, error) {
	var user user.User
	err := r.getDB(ctx).WithContext(ctx).Where("id = ?", id).First(&user).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) FindByEmail(ctx context.Context, email string) (*user.User, error) {
	var user user.User
	err := r.getDB(ctx).WithContext(ctx).Where("email = ?", email).First(&user).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

// WithTransaction executes a function within a database transaction
// If the function returns an error, the transaction is rolled back
func (r *UserRepository) WithTransaction(ctx context.Context, fn func(*gorm.DB) error) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		return fn(tx)
	})
}

// GetDB returns the underlying GORM database instance for advanced operations
func (r *UserRepository) GetDB() *gorm.DB {
	return r.db
}
