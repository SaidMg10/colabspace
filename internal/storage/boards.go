package storage

import (
	"context"
	"database/sql"

	"github.com/SaidMg10/colabspace/internal/app/board"
	"github.com/SaidMg10/colabspace/internal/config"
)

type BoardStore struct {
	db *sql.DB
}

func (s *BoardStore) Create(ctx context.Context, brd *board.Board) error {
	query := `
		INSERT INTO boards (name, description, user_id)
		VALUES ($1, $2, $3)
		RETURNING id, created_at
	`
	ctx, cancel := context.WithTimeout(ctx, config.QueryTimeoutDuration)
	defer cancel()

	err := s.db.QueryRowContext(
		ctx,
		query,
		brd.Name,
		brd.Description,
		brd.UserID,
	).Scan(
		&brd.ID,
		&brd.CreatedAt,
	)
	if err != nil {
		return err
	}
	return nil
}

func (s *BoardStore) Get(ctx context.Context) ([]board.Board, error) {
	query := `
		SELECT id, name, description, user_id, created_at, updated_at, deleted_at
		FROM boards
	`
	ctx, cancel := context.WithTimeout(ctx, config.QueryTimeoutDuration)
	defer cancel()

	rows, err := s.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var boards []board.Board

	for rows.Next() {
		var b board.Board

		err := rows.Scan(
			&b.ID,
			&b.Name,
			&b.Description,
			&b.UserID,
			&b.CreatedAt,
			&b.UpdatedAt,
			&b.DeletedAt,
		)
		if err != nil {
			return nil, err
		}
		boards = append(boards, b)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return boards, nil
}

func (s *BoardStore) GetBoardsForUser(ctx context.Context, userID int64) ([]board.Board, error) {
	query := `
		SELECT b.id, b.name, b.description
		FROM boards b
		LEFT JOIN board_users bu ON b.id = bu.board_id
		WHERE b.user_id = $1 OR bu.user_id = $1
		GROUP BY b.id
	`
	rows, err := s.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var boards []board.Board
	for rows.Next() {
		var b board.Board
		err := rows.Scan(&b.ID, &b.Name, &b.Description)
		if err != nil {
			return nil, err
		}
		boards = append(boards, b)
	}
	return boards, nil
}

func (s *BoardStore) GetById(ctx context.Context, id int64) (*board.Board, error) {
	query := `
		SELECT id, name, description, user_id, created_at, updated_at, deleted_at
		FROM boards
		WHERE id = $1
	`

	ctx, cancel := context.WithTimeout(ctx, config.QueryTimeoutDuration)
	defer cancel()

	row := s.db.QueryRowContext(ctx, query, id)

	var brd board.Board

	err := row.Scan(
		&brd.ID,
		&brd.Name,
		&brd.Description,
		&brd.UserID,
		&brd.CreatedAt,
		&brd.UpdatedAt,
		&brd.DeletedAt,
	)
	if err != nil {
		return nil, err
	}

	return &brd, nil
}

func (s *BoardStore) Update(ctx context.Context, brd *board.Board) error {
	query := `
		UPDATE boards
		SET name = $1, description = $2, updated_at = NOW()
		WHERE id = $3
		RETURNING id, updated_at
	`

	ctx, cancel := context.WithTimeout(ctx, config.QueryTimeoutDuration)
	defer cancel()

	err := s.db.QueryRowContext(
		ctx,
		query,
		brd.Name,
		brd.Description,
		brd.ID,
	).Scan(
		&brd.ID,
		&brd.UpdatedAt,
	)
	if err != nil {
		return err
	}
	return nil
}

func (s *BoardStore) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM boards WHERE id = $1`

	ctx, cancel := context.WithTimeout(ctx, config.QueryTimeoutDuration)
	defer cancel()

	_, err := s.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	return nil
}
