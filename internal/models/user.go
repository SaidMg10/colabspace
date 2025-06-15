package models

type User struct {
	ID        int64    `json:"id"`
	Username  string   `json:"username"`
	FirstName string   `json:"first_name"`
	LastName  string   `json:"last_name"`
	Email     string   `json:"email"`
	Password  Password `json:"-"`
	IsActive  bool     `json:"is_active"`
	CreatedAt string   `json:"created_at"`
	UpdatedAt string   `json:"updated_at"`
	DeletedAt string   `json:"deleted_at"`
	Role      Role     `json:"role"`
}
