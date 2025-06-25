package boardusers

import (
	"context"

	"github.com/SaidMg10/colabspace/internal/app/user"
)

type BoardUsersRepository interface {
	Add(ctx context.Context, boardUsers *BoardUsers) error
	GetMembers(ctx context.Context, boardID int64) ([]user.User, error)
	Leave(ctx context.Context, boardID int64, userID int64) error
	Delete(ctx context.Context, boardID int64, userID int64) error
	IsMember(ctx context.Context, boardID, userID int64) (bool, error)
}
