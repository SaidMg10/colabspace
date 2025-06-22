package board

import "context"

type BoardRepository interface {
	Create(ctx context.Context, board *Board) error
	Get(ctx context.Context) ([]Board, error)
	GetById(ctx context.Context, id int64) (*Board, error)
	Update(ctx context.Context, board *Board) error
	Delete(ctx context.Context, id int64) error
}
