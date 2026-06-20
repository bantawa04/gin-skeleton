package user

// UserCreateRequest represents the request payload for creating a new user
type UserCreateRequest struct {
	Name     string `json:"name" binding:"required,min=2,max=100"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6,max=100"`
}

// UserUpdateRequest represents the request payload for updating an existing user
type UserUpdateRequest struct {
	Name     *string `json:"name,omitempty" binding:"omitempty,min=2,max=100"`
	Email    *string `json:"email,omitempty" binding:"omitempty,email"`
	Password *string `json:"password,omitempty" binding:"omitempty,min=6,max=100"`
}

// SignupInput represents data needed to create a user during signup
type SignupInput struct {
	FirstName string
	LastName  string
	Email     string
}
