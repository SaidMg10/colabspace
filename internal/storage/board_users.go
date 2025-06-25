package storage

import (
	"context"
	"database/sql"
	"fmt"

	boardusers "github.com/SaidMg10/colabspace/internal/app/board_users"
	"github.com/SaidMg10/colabspace/internal/app/user"
	"github.com/SaidMg10/colabspace/internal/config"
)

type BoardUsersStore struct {
	db *sql.DB
}

func (s *BoardUsersStore) Add(ctx context.Context, boardUsers *boardusers.BoardUsers) error {
	// query
	query := `
		INSERT INTO board_users (board_id, user_id)
		VALUES ($1, $2)
	`
	ctx, cancel := context.WithTimeout(ctx, config.QueryTimeoutDuration)
	defer cancel()
	// rowContext
	_, err := s.db.ExecContext(
		ctx,
		query,
		boardUsers.BoardID,
		boardUsers.UserID,
	)
	// if error
	if err != nil {
		return fmt.Errorf("failed to insert board user: %w", err)
	}

	return nil
}

func (s *BoardUsersStore) GetMembers(ctx context.Context, boardID int64) ([]user.User, error) {
	query := `
		SELECT u.id, u.username, u.first_name, u.last_name
		FROM users u
		JOIN board_users bu ON u.id = bu.user_id
		WHERE bu.board_id = $1
	`

	ctx, cancel := context.WithTimeout(ctx, config.QueryTimeoutDuration)
	defer cancel()

	rows, err := s.db.QueryContext(ctx, query, boardID)
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

func (s *BoardUsersStore) Leave(ctx context.Context, boardID, userID int64) error {
	query := `
	DELETE FROM board_users WHERE board_id = $1 AND user_id = $2
	`
	ctx, cancel := context.WithTimeout(ctx, config.QueryTimeoutDuration)
	defer cancel()

	_, err := s.db.ExecContext(
		ctx,
		query,
		boardID,
		userID,
	)
	if err != nil {
		return err
	}

	return nil
}

func (s *BoardUsersStore) Delete(ctx context.Context, boardID, userID int64) error {
	query := `
	DELETE FROM board_users WHERE board_id = $1 AND user_id = $2
	`
	ctx, cancel := context.WithTimeout(ctx, config.QueryTimeoutDuration)
	defer cancel()

	_, err := s.db.ExecContext(
		ctx,
		query,
		boardID,
		userID,
	)
	if err != nil {
		return err
	}

	return nil
}

func (s *BoardUsersStore) IsMember(ctx context.Context, boardID, userID int64) (bool, error) {
	query := `
		SELECT EXISTS (
		SELECT 1 FROM board_users
		WHERE board_id = $1 AND user_id = $2
	)`
	var exists bool
	err := s.db.QueryRowContext(
		ctx,
		query,
		boardID,
		userID,
	).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}
