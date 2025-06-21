package user

type CreateUserRequest struct {
	Username  string `json:"username" validate:"required,max=255"`
	FirstName string `json:"first_name" validate:"required,max=255"`
	LastName  string `json:"last_name" validate:"required,max=255"`
	Email     string `json:"email" validate:"required,email,max=255"`
	Password  string `json:"password" validate:"required,min=6,max=100"`
}

type UpdateUserRequest struct {
	FirstName *string `json:"first_name" validate:"omitempty,max=255"`
	LastName  *string `json:"last_name" validate:"omitempty,max=255"`
	Password  *string `json:"password" validate:"omitempty,min=6,max=100"`
}
