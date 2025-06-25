package boardusers

type BoardUsers struct {
	ID      int64 `json:"id"`
	BoardID int64 `json:"board_id"`
	UserID  int64 `json:"user_id"`
}
