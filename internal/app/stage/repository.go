package stage

import "context"

type StageRepository interface {
	Create(ctx context.Context, stage *Stage) error
	Get(ctx context.Context, boardID int64) ([]Stage, error)
	GetById(ctx context.Context, boardID, stageID int64) (*Stage, error)
	Update(ctx context.Context, stage *Stage) error
	Delete(ctx context.Context, boardID, stageID int64) error
}
