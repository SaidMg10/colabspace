package board

func MapCreateBoardRequestToBoard(req CreateBoardRequest) (*Board, error) {
	board := &Board{
		Name:        req.Name,
		Description: req.Description,
		UserID:      req.UserID,
	}
	return board, nil
}

func MapUpdateBoardRequestToBoard(req UpdateBoardRequest, boardExist *Board) (*Board, error) {
	updated := *boardExist
	if req.Name != nil {
		updated.Name = *req.Name
	}
	if req.Description != nil {
		updated.Description = *req.Description
	}
	return &updated, nil
}

func MapToBoardsResponse(b []Board) []BoardResponse {
	resp := make([]BoardResponse, len(b))
	for i, v := range b {
		resp[i] = BoardResponse{
			ID:          v.ID,
			Name:        v.Name,
			Description: v.Description,
			UserID:      v.UserID,
		}
	}
	return resp
}

func MapToBoardsForUsersResponse(b []Board) []BoardForUsersResponse {
	resp := make([]BoardForUsersResponse, len(b))
	for i, v := range b {
		resp[i] = BoardForUsersResponse{
			ID:          v.ID,
			Name:        v.Name,
			Description: v.Description,
		}
	}
	return resp
}

func MapToBoardResponse(b *Board) *BoardResponse {
	resp := &BoardResponse{
		ID:          b.ID,
		Name:        b.Name,
		Description: b.Description,
		UserID:      b.UserID,
	}

	return resp
}
