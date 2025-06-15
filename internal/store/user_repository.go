package store

import (
	"context"

	"github.com/SaidMg10/colabspace/internal/models"
)

type UserRepository interface {
	// TODO: Create agregar tx en el momento de agregar invitacion
	Create(ctx context.Context, user *models.User) error
	Get(ctx context.Context) ([]models.User, error)
	GetByID(ctx context.Context, id int64) (*models.User, error)
	GetByEmail(ctx context.Context, email string) (*models.User, error)
	Update(ctx context.Context, user *models.User) error
	Delete(ctx context.Context, id int64) error
}
