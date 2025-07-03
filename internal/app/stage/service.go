package stage

import (
	"context"

	"github.com/SaidMg10/colabspace/internal/utils"
	"go.uber.org/zap"
)

type Service struct {
	// Se inyecta la interfaz del repositorio (Repository),
	// permitiendo que el Service use cualquier implementación concreta (por ejemplo, UserStore).
	Repo   StageRepository
	Logger *zap.SugaredLogger
}

// NewService es el constructor del Service.
// Recibe una implementación del repositorio y retorna una instancia del Service.
// Esto permite desacoplar el Service de la implementación concreta de almacenamiento.
func NewService(repo StageRepository, logger *zap.SugaredLogger) *Service {
	return &Service{
		Repo:   repo,
		Logger: logger,
	}
}

func (s *Service) Create(ctx context.Context, cSR *CreateStageRequest) (*CreateStageResponse, error) {
	stage, err := MapCreateStageRequestToStage(*cSR)
	if err != nil {
		return nil, err
	}

	if err := s.Repo.Create(ctx, stage); err != nil {
		return nil, err
	}

	resp := &CreateStageResponse{
		ID:        stage.ID,
		CreatedAt: stage.CreatedAt,
	}
	return resp, nil
}

func (s *Service) Get(ctx context.Context, boardID int64) ([]StagePublicResponse, error) {
	stages, err := s.Repo.Get(ctx, boardID)
	if err != nil {
		return nil, err
	}
	resp := MapToStagesPublicResponse(stages)

	return resp, nil
}

func (s *Service) GetById(ctx context.Context, boardID, stageID int64) (*StagePublicResponse, error) {
	stage, err := s.Repo.GetById(ctx, boardID, stageID)
	if err != nil {
		return nil, err
	}

	if stage == nil {
		return nil, utils.ErrNotFound
	}

	resp := MapToStagePublicResponse(stage)

	return resp, nil
}

func (s *Service) Update(ctx context.Context, boardID, stageID int64, uSR *UpdateStageRequest) (*UpdateStageResponse, error) {
	// Obtenemos el stage
	stageReq, err := s.Repo.GetById(ctx, boardID, stageID)
	if err != nil {
		return nil, err
	}
	// Verificamos que exista
	if stageReq == nil {
		return nil, utils.ErrNotFound
	}
	// Validamos y Pasamos a una funcion todos los datos a validar
	stage, err := MapUpdateStageRequestToStage(*uSR, stageReq)
	if err != nil {
		return nil, err
	}
	if err := s.Repo.Update(ctx, stage); err != nil {
		return nil, err
	}
	resp := &UpdateStageResponse{
		ID:        stage.ID,
		UpdatedAt: stage.UpdatedAt,
	}
	return resp, nil
}

func (s *Service) Delete(ctx context.Context, boardID, stageID int64) error {
	if err := s.Repo.Delete(ctx, boardID, stageID); err != nil {
		return err
	}
	return nil
}
