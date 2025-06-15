package store

import (
	"context"
	"database/sql"
	"errors"

	"github.com/SaidMg10/colabspace/internal/models"
)

var (
	ErrDuplicateEmail    = errors.New("a user with that email already exists")
	ErrDuplicateUsername = errors.New("a user with that username already exists")
)

type UserStore struct {
	db *sql.DB
}

func (s *UserStore) Create(ctx context.Context, user *models.User) error {
	query := `
	INSERT INTO users (username, first_name, last_name, email, password, role)
	VALUES ($1, $2, $3, $4, $5, $6)
	RETURNING id, created_at
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	err := s.db.QueryRowContext(
		ctx,
		query,
		user.Username,
		user.FirstName,
		user.LastName,
		user.Email,
		user.Password.Hash(),
		int(user.Role),
	).Scan(
		&user.ID,
		&user.CreatedAt,
	)
	if err != nil {
		switch {
		case err.Error() == `pq: duplicate key value violates unique constraint "users_email_key"`:
			return ErrDuplicateEmail
		case err.Error() == `pq: duplicate key value violates unique constraint "users_username_key"`:
			return ErrDuplicateUsername
		default:
			return err
		}
	}
	return nil
}

func (s *UserStore) Get(ctx context.Context) ([]models.User, error) {
	return nil, nil
}

func (s *UserStore) GetByID(ctx context.Context, id int64) (*models.User, error) {
	return nil, nil
}

func (s *UserStore) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	return nil, nil
}

func (s *UserStore) Update(ctx context.Context, user *models.User) error {
	return nil
}

func (s *UserStore) Delete(ctx context.Context, id int64) error {
	return nil
}
