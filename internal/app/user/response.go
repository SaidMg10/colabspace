package user

import (
	"time"
)

type CreateUserResponse struct {
	ID        int64  `json:"id"`
	CreatedAt string `json:"created_at"`
}

type UpdateUserResponse struct {
	ID        int64  `json:"id"`
	UpdatedAt string `json:"updated_at"`
}

type UserResponse struct {
	ID        int64      `json:"id"`
	Username  string     `json:"username"`
	FirstName string     `json:"first_name"`
	LastName  string     `json:"last_name"`
	Email     string     `json:"email"`
	IsActive  bool       `json:"is_active"`
	Role      Role       `json:"role"`
	CreatedAt string     `json:"created_at"`
	UpdatedAt string     `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at"`
}

type UserSummaryResponse struct {
	ID        int64  `json:"id"`
	Username  string `json:"username"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}
