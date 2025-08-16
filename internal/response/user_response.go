package response

import (
	"gin/internal/models"
)

// UserResponse represents the user data sent to the client
type UserResponse struct {
	ID               string  `json:"id"`
	FirstName        *string `json:"first_name,omitempty"`
	LastName         *string `json:"last_name,omitempty"`
	Email            string  `json:"email"`
	Phone            *string `json:"phone,omitempty"`
	Province         *string `json:"province,omitempty"`
	District         *string `json:"district,omitempty"`
	Address          *string `json:"address,omitempty"`
	Type             string  `json:"type"`
	SocialProvider   *string `json:"social_provider,omitempty"`
	SocialProviderID *string `json:"social_provider_id,omitempty"`
	LastSignInAt     *string `json:"last_sign_in_at,omitempty"`
	Status           string  `json:"status"`
	CreatedAt        string  `json:"created_at"`
	UpdatedAt        string  `json:"updated_at"`
	FullName         string  `json:"full_name,omitempty"`
}

// UserListResponse represents a list of users
type UserListResponse struct {
	Users []*UserResponse `json:"users"`
	Total int64           `json:"total"`
}

// PaginatedUserResponse represents a paginated list of users
type PaginatedUserResponse struct {
	Users []*UserResponse `json:"users"`
	Meta  PaginationMeta  `json:"meta"`
}

// PaginationMeta represents pagination metadata
type PaginationMeta struct {
	Page       int   `json:"page"`
	TotalPages int   `json:"totalPages"`
	PerPage    int   `json:"perPage"`
	TotalItems int64 `json:"totalItems"`
}

// ToUserResponse converts a User model to UserResponse DTO
func ToUserResponse(user *models.User) *UserResponse {
	if user == nil {
		return nil
	}

	response := &UserResponse{
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
		Status:           string(user.Status),
		CreatedAt:        user.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:        user.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
		FullName:         user.FullName(),
	}

	// Handle LastSignInAt if it exists
	if user.LastSignInAt != nil {
		formattedTime := user.LastSignInAt.Format("2006-01-02T15:04:05Z07:00")
		response.LastSignInAt = &formattedTime
	}

	return response
}

// ToPaginatedUserResponse converts a slice of User models to PaginatedUserResponse DTO
func ToPaginatedUserResponse(users []*models.User, page, perPage int, totalItems int64) *PaginatedUserResponse {
	if users == nil {
		return &PaginatedUserResponse{
			Users: []*UserResponse{},
			Meta: PaginationMeta{
				Page:       page,
				TotalPages: 0,
				PerPage:    perPage,
				TotalItems: 0,
			},
		}
	}

	userResponses := make([]*UserResponse, len(users))
	for i, user := range users {
		userResponses[i] = ToUserResponse(user)
	}

	// Calculate total pages
	totalPages := int((totalItems + int64(perPage) - 1) / int64(perPage))

	return &PaginatedUserResponse{
		Users: userResponses,
		Meta: PaginationMeta{
			Page:       page,
			TotalPages: totalPages,
			PerPage:    perPage,
			TotalItems: totalItems,
		},
	}
}
