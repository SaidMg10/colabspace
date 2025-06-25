package task

import "time"

type Task struct {
	ID          int64
	Title       string
	Description string
	UserIDs     []int64
	Completed   bool
	ExpiredAt   time.Time
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
