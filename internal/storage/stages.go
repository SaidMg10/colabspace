package storage

import (
	"context"
	"database/sql"

	"github.com/SaidMg10/colabspace/internal/app/stage"
)

type StageStore struct {
	db *sql.DB
}

func (s *StageStore) Create(ctx context.Context, stage *stage.Stage) error {
	return nil
}

func (s *StageStore) Get(ctx context.Context) ([]stage.Stage, error) {
	return nil, nil
}

func (s *StageStore) GetById(ctx context.Context, stageID int64) (*stage.Stage, error) {
	return nil, nil
}

func (s *StageStore) Update(ctx context.Context, state *stage.Stage) error {
	return nil
}

func (s *StageStore) Delete(ctx context.Context, stageID int64) error {
	return nil
}
