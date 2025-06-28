package stage

import "context"

type StageRepository interface {
	Create(ctx context.Context, stage *Stage) error
	Get(ctx context.Context) ([]Stage, error)
	GetById(ctx context.Context, stageID int64) (*Stage, error)
	Update(ctx context.Context, stage *Stage) error
	Delete(ctx context.Context, stageID int64) error
}
