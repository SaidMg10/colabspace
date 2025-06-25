package boardusers

type CreateBoardUsersRequest struct {
	BoardID int64  `json:"-"`
	Value   string `json:"value" validate:"required"`
}
