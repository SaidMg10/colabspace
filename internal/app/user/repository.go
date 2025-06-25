package user

import "context"

type UserRepository interface {
	Create(ctx context.Context, user *User) error
	Get(ctx context.Context) ([]User, error)
	GetById(ctx context.Context, id int64) (*User, error)
	GetByEmail(ctx context.Context, email string) (*User, error)
	GetByUsernameOrEmail(ctx context.Context, value string) (*User, error)
	Update(ctx context.Context, user *User) error
	Delete(ctx context.Context, id int64) error
}
