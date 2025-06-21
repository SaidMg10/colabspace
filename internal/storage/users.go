package storage

import (
	"context"
	"database/sql"
	"errors"

	"github.com/SaidMg10/colabspace/internal/app/user"
	"github.com/SaidMg10/colabspace/internal/config"
	utils "github.com/SaidMg10/colabspace/internal/utils"
	"github.com/jackc/pgx/v5/pgconn"
)

type UserStore struct {
	db *sql.DB
}

func (s *UserStore) Create(ctx context.Context, u *user.User) error {
	query := `
	INSERT INTO users (username, first_name, last_name, email, password, role)
	VALUES ($1, $2, $3, $4, $5, $6)
	RETURNING id, created_at
	`

	ctx, cancel := context.WithTimeout(ctx, config.QueryTimeoutDuration)
	defer cancel()

	err := s.db.QueryRowContext(
		ctx,
		query,
		u.Username,
		u.FirstName,
		u.LastName,
		u.Email,
		u.Password.Hash,
		int(u.Role),
	).Scan(
		&u.ID,
		&u.CreatedAt,
	)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == "23505" { // Unique violation
				switch pgErr.ConstraintName {
				case "users_email_key":
					return user.ErrDuplicateEmail
				case "users_username_key":
					return user.ErrDuplicateUsername
				}
			}
		}
		return err
	}

	return nil
}

func (s *UserStore) Get(ctx context.Context) ([]user.User, error) {
	query := `
	SELECT id, username, first_name, last_name,
	email, role, created_at, updated_at, deleted_at
	FROM users
	`

	ctx, cancel := context.WithTimeout(ctx, config.QueryTimeoutDuration)
	defer cancel()

	rows, err := s.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []user.User

	for rows.Next() {
		var u user.User

		err := rows.Scan(
			&u.ID,
			&u.Username,
			&u.FirstName,
			&u.LastName,
			&u.Email,
			&u.Role,
			&u.CreatedAt,
			&u.UpdatedAt,
			&u.DeletedAt,
		)
		if err != nil {
			return nil, err
		}

		users = append(users, u)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

func (s *UserStore) GetById(ctx context.Context, id int64) (*user.User, error) {
	query := `
	SELECT id, username, first_name, last_name,
	email, role, created_at, updated_at, deleted_at
	FROM users
	WHERE id = $1
	`

	ctx, cancel := context.WithTimeout(ctx, config.QueryTimeoutDuration)
	defer cancel()

	row := s.db.QueryRowContext(ctx, query, id)

	var u user.User

	err := row.Scan(
		&u.ID,
		&u.Username,
		&u.FirstName,
		&u.LastName,
		&u.Email,
		&u.Role,
		&u.CreatedAt,
		&u.UpdatedAt,
		&u.DeletedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, utils.ErrNotFound
		}
		return nil, err
	}

	return &u, nil
}

func (s *UserStore) GetByEmail(ctx context.Context, email string) (*user.User, error) {
	query := `
		SELECT id, username, email, password, created_at FROM users
		WHERE email = $1 AND is_active = true
	`

	ctx, cancel := context.WithTimeout(ctx, config.QueryTimeoutDuration)
	defer cancel()

	user := &user.User{}
	err := s.db.QueryRowContext(ctx, query, email).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.Password.Hash,
		&user.CreatedAt,
	)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return nil, utils.ErrNotFound
		default:
			return nil, err
		}
	}

	return user, nil
}

func (s *UserStore) Update(ctx context.Context, u *user.User) error {
	query := `
		UPDATE users
	  SET first_name = $1, last_name = $2, password = COALESCE(NULLIF($3, ''), password), updated_at = NOW()
		WHERE id = $4
		RETURNING id, updated_at
	`

	ctx, cancel := context.WithTimeout(ctx, config.QueryTimeoutDuration)
	defer cancel()

	err := s.db.QueryRowContext(
		ctx,
		query,
		u.FirstName,
		u.LastName,
		u.Password.Hash,
		u.ID,
	).Scan(
		&u.ID,
		&u.UpdatedAt,
	)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return utils.ErrNotFound
		default:
			return err
		}
	}

	return nil
}

func (s *UserStore) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM users WHERE id = $1`

	ctx, cancel := context.WithTimeout(ctx, config.QueryTimeoutDuration)
	defer cancel()

	_, err := s.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	return nil
}
