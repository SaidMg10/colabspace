package storage

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/SaidMg10/colabspace/internal/app/stage"
	"github.com/SaidMg10/colabspace/internal/config"
)

type StageStore struct {
	db *sql.DB
}

func (s *StageStore) Create(ctx context.Context, stage *stage.Stage) error {
	query := `
		INSERT INTO stages (title, position, wip_limit, stage_type, board_id)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, created_at;
	`

	ctx, cancel := context.WithTimeout(ctx, config.QueryTimeoutDuration)
	defer cancel()

	err := s.db.QueryRowContext(
		ctx,
		query,
		stage.Title,
		stage.Position,
		stage.WipLimit,
		stage.StageType,
		stage.BoardId,
	).Scan(
		&stage.ID,
		&stage.CreatedAt,
	)
	if err != nil {
		return err
	}

	return nil
}

func (s *StageStore) Get(ctx context.Context, boardID int64) ([]stage.Stage, error) {
	query := `	
    SELECT id, title, position, wip_limit, stage_type, board_id, created_at, updated_at
    FROM stages
    WHERE board_id = $1;
	`

	ctx, cancel := context.WithTimeout(ctx, config.QueryTimeoutDuration)
	defer cancel()

	rows, err := s.db.QueryContext(
		ctx,
		query,
		boardID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var stages []stage.Stage

	for rows.Next() {
		var st stage.Stage
		err := rows.Scan(
			&st.ID,
			&st.Title,
			&st.Position,
			&st.WipLimit,
			&st.StageType,
			&st.BoardId,
			&st.CreatedAt,
			&st.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed scanning stage row: %w", err)
		}
		stages = append(stages, st)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return stages, nil
}

func (s *StageStore) GetById(ctx context.Context, boardID, stageID int64) (*stage.Stage, error) {
	query := `	
    SELECT id, title, position, wip_limit, stage_type, board_id, created_at, updated_at
    FROM stages
    WHERE id = $1 AND board_id = $2;
	`
	ctx, cancel := context.WithTimeout(ctx, config.QueryTimeoutDuration)
	defer cancel()

	var st stage.Stage
	err := s.db.QueryRowContext(
		ctx,
		query,
		stageID,
		boardID,
	).Scan(
		&st.ID,
		&st.Title,
		&st.Position,
		&st.WipLimit,
		&st.StageType,
		&st.BoardId,
		&st.CreatedAt,
		&st.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &st, nil
}

func (s *StageStore) Update(ctx context.Context, stage *stage.Stage) error {
	query := `
		UPDATE stages
    SET title = $1, position = $2, wip_limit = $3, stage_type = $4, updated_at = NOW()
    WHERE id = $5 AND board_id = $6
    RETURNING id, updated_at;
	`
	ctx, cancel := context.WithTimeout(ctx, config.QueryTimeoutDuration)
	defer cancel()

	err := s.db.QueryRowContext(
		ctx,
		query,
		stage.Title,
		stage.Position,
		stage.WipLimit,
		stage.StageType,
		stage.ID,
		stage.BoardId,
	).Scan(
		&stage.ID,
		&stage.UpdatedAt,
	)
	if err != nil {
		return err
	}

	return nil
}

func (s *StageStore) Delete(ctx context.Context, boardID, stageID int64) error {
	query := `
		DELETE FROM stages
    WHERE id = $1 AND board_id = $2;
	`
	ctx, cancel := context.WithTimeout(ctx, config.QueryTimeoutDuration)
	defer cancel()

	_, err := s.db.ExecContext(ctx, query, stageID, boardID)
	if err != nil {
		return err
	}

	return nil
}
