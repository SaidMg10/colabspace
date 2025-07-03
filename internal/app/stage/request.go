package stage

import "time"

type CreateStageRequest struct {
	ID        int64     `json:"id"`
	Title     string    `json:"title"`
	Position  int       `json:"position"`
	WipLimit  int       `json:"wip_limit"`
	StageType int       `json:"stage_type"`
	BoardId   int64     `json:"-"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type UpdateStageRequest struct {
	Title     *string `json:"title" validate:"omitempty"`
	Position  *int    `json:"position" validate:"omitempty"`
	WipLimit  *int    `json:"wip_limit" validate:"omitempty"`
	StageType *int    `json:"stage_type" validate:"omitempty"`
}
