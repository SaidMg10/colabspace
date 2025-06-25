package boardusers

func MapCreateBoardUsersRequestToBoardUsers(boardID, userID int64) (*BoardUsers, error) {
	boardUsers := &BoardUsers{
		BoardID: boardID,
		UserID:  userID,
	}
	return boardUsers, nil
}
