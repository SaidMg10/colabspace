package storage

import (
	"context"
	"database/sql"

	"github.com/SaidMg10/colabspace/internal/app/board"
	"github.com/SaidMg10/colabspace/internal/app/user"
)

type Storage struct {
	Users  user.UserRepository
	Boards board.BoardRepository
}

func NewStorage(db *sql.DB) Storage {
	return Storage{
		Users:  &UserStore{db},
		Boards: &BoardStore{db},
	}
}

func withTx(db *sql.DB, ctx context.Context, fn func(*sql.Tx) error) error {
	tx, err := db.BeginTx(ctx, nil) // 1. Inicia transacción
	if err != nil {
		return err
	}

	if err := fn(tx); err != nil { // 2. Ejecuta callback
		_ = tx.Rollback() // 3. Rollback si falla
		return err
	}

	return tx.Commit() // 4. Commit si todo va bien
}
