package board

type CreateBoardRequest struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	Description string `json:"Description"`
	UserID      int64  `json:"user_id" validate:"omitempty,max=255"`
}

type UpdateBoardRequest struct {
	Name        *string `json:"name"`
	Description *string `json:"description"`
	UpdatedAt   *string `json:"updated_at"`
}
