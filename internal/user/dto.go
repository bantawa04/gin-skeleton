package user

import (
	"time"

	"gin/internal/utils/transformer"
)

// UserDTO represents the data transfer object for User
type UserDTO struct {
	ID               string     `json:"id"`
	FirstName        *string    `json:"firstName,omitempty"`
	LastName         *string    `json:"lastName,omitempty"`
	Email            string     `json:"email"`
	Phone            *string    `json:"phone,omitempty"`
	Province         *string    `json:"province,omitempty"`
	District         *string    `json:"district,omitempty"`
	Address          *string    `json:"address,omitempty"`
	Type             string     `json:"type"`
	SocialProvider   *string    `json:"socialProvider,omitempty"`
	SocialProviderID *string    `json:"socialProviderId,omitempty"`
	LastSignInAt     *time.Time `json:"lastSignInAt,omitempty"`
	Status           string     `json:"status"`
	CreatedAt        time.Time  `json:"createdAt"`
	UpdatedAt        time.Time  `json:"updatedAt"`
	FullName         string     `json:"fullName,omitempty"`
}

// PaginationMeta represents pagination metadata
type PaginationMeta struct {
	Page       int   `json:"page"`
	TotalPages int   `json:"totalPages"`
	PerPage    int   `json:"perPage"`
	TotalItems int64 `json:"totalItems"`
}

// PaginatedUserDTO represents a paginated list of users
type PaginatedUserDTO struct {
	Users []UserDTO      `json:"users"`
	Meta  PaginationMeta `json:"meta"`
}

// FromUserModel converts a User model to a UserDTO
func FromUserModel(user User) UserDTO {
	dto := UserDTO{
		ID:               user.ID,
		FirstName:        user.FirstName,
		LastName:         user.LastName,
		Email:            user.Email,
		Phone:            user.Phone,
		Province:         user.Province,
		District:         user.District,
		Address:          user.Address,
		Type:             string(user.Type),
		SocialProvider:   user.SocialProvider,
		SocialProviderID: user.SocialProviderID,
		LastSignInAt:     user.LastSignInAt,
		Status:           string(user.Status),
		CreatedAt:        user.CreatedAt,
		UpdatedAt:        user.UpdatedAt,
		FullName:         user.FullName(),
	}

	return dto
}

// TransformUserCollection transforms a slice of User models to a slice of UserDTOs
func TransformUserCollection(users []*User) []UserDTO {
	// Convert []*User to []User for the transformer
	modelSlice := make([]User, len(users))
	for i, user := range users {
		if user != nil {
			modelSlice[i] = *user
		}
	}
	return transformer.TransformCollection(modelSlice, FromUserModel)
}

// ToPaginatedUserDTO creates a paginated user DTO
func ToPaginatedUserDTO(users []*User, page, perPage int, totalItems int64) PaginatedUserDTO {
	userDTOs := TransformUserCollection(users)

	// Calculate total pages
	totalPages := int((totalItems + int64(perPage) - 1) / int64(perPage))

	return PaginatedUserDTO{
		Users: userDTOs,
		Meta: PaginationMeta{
			Page:       page,
			TotalPages: totalPages,
			PerPage:    perPage,
			TotalItems: totalItems,
		},
	}
}
