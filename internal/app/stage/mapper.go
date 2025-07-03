package stage

func MapCreateStageRequestToStage(req CreateStageRequest) (*Stage, error) {
	stage := &Stage{
		Title:     req.Title,
		Position:  req.Position,
		WipLimit:  req.WipLimit,
		StageType: req.StageType,
		BoardId:   req.BoardId,
	}

	return stage, nil
}

func MapToStagesPublicResponse(s []Stage) []StagePublicResponse {
	resp := make([]StagePublicResponse, len(s))
	for i, v := range s {
		resp[i] = StagePublicResponse{
			ID:        v.ID,
			Title:     v.Title,
			Position:  v.Position,
			WipLimit:  v.WipLimit,
			StageType: v.StageType,
		}
	}
	return resp
}

func MapToStagesPrivateResponse(s []Stage) []StagePrivateResponse {
	resp := make([]StagePrivateResponse, len(s))
	for i, v := range s {
		resp[i] = StagePrivateResponse(v)
	}
	return resp
}

func MapToStagePublicResponse(s *Stage) *StagePublicResponse {
	resp := &StagePublicResponse{
		ID:        s.ID,
		Title:     s.Title,
		Position:  s.Position,
		WipLimit:  s.WipLimit,
		StageType: s.StageType,
	}
	return resp
}

func MapUpdateStageRequestToStage(req UpdateStageRequest, stageExist *Stage) (*Stage, error) {
	updated := *stageExist
	if req.Title != nil {
		updated.Title = *req.Title
	}
	if req.Position != nil {
		updated.Position = *req.Position
	}
	if req.WipLimit != nil {
		updated.WipLimit = *req.WipLimit
	}
	if req.StageType != nil {
		updated.StageType = *req.StageType
	}

	// Devolvemos el Stage preparado para ser enviado a la bd
	return &updated, nil
}
