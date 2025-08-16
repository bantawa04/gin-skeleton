package request

type UserCreateRequest struct {
	FirstName string `json:"first_name" validate:"required,min=2,max=255"`
	LastName  string `json:"last_name" validate:"required,min=2,max=255"`
	Email     string `json:"email" validate:"required,email,min=2,max=255"`
	Password  string `json:"password" validate:"required,min=8,max=255"`
}

type UserUpdateRequest struct {
	FirstName string `json:"first_name" validate:"required,min=2,max=255"`
	LastName  string `json:"last_name" validate:"required,min=2,max=255"`
	Email     string `json:"email" validate:"required,email,min=2,max=255"`
	Password  string `json:"password,omitempty" validate:"omitempty,min=8,max=255"`
}
