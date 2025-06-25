package board

import "time"

type CreateBoardResponse struct {
	ID        int64     `json:"id"`
	CreatedAt time.Time `json:"created_at"`
}

type UpdateBoardResponse struct {
	ID        int64     `json:"id"`
	UpdatedAt time.Time `json:"updated_at"`
}

type BoardResponse struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	Description string `json:"Description"`
	UserID      int64  `json:"user_id"`
}

type BoardForUsersResponse struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	Description string `json:"Description"`
}
