package dto

type CreateUserDto struct {
	Username  string `json:"username" validate:"required,max=255"`
	FirstName string `json:"first_name" validate:"required,max=255"`
	LastName  string `json:"last_name" validate:"required,max=255"`
	Email     string `json:"email" validate:"required,email,max=255"`
	Password  string `json:"password" validate:"required,min=6,max=100"`
}
type CreateUserResponse struct {
	ID        int64  `json:"id"`
	CreatedAt string `json:"created_at"`
}
