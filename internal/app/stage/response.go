package stage

import "time"

type CreateStageResponse struct {
	ID        int64     `json:"id"`
	CreatedAt time.Time `json:"created_at"`
}

type UpdateStageResponse struct {
	ID        int64     `json:"id"`
	UpdatedAt time.Time `json:"updated_at"`
}

type StagePrivateResponse struct {
	ID        int64     `json:"id"`
	Title     string    `json:"title"`
	Position  int       `json:"position"`
	WipLimit  int       `json:"wip_limit"`
	StageType int       `json:"stage_type"`
	BoardId   int64     `json:"board_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type StagePublicResponse struct {
	ID        int64  `json:"id"`
	Title     string `json:"title"`
	Position  int    `json:"position"`
	WipLimit  int    `json:"wip_limit"`
	StageType int    `json:"stage_type"`
}
