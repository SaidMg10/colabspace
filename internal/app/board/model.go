package board

import "time"

type Board struct {
	ID          int64
	Name        string
	Description string
	UserID      int64
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   *time.Time
}
