package board

import (
	"context"

	"github.com/SaidMg10/colabspace/internal/utils"
	"go.uber.org/zap"
)

type Service struct {
	// Se inyecta la interfaz del repositorio (UserRepository),
	// permitiendo que el Service use cualquier implementación concreta (por ejemplo, UserStore).
	Repo   BoardRepository
	Logger *zap.SugaredLogger
}

// NewService es el constructor del Service.
// Recibe una implementación del repositorio y retorna una instancia del Service.
// Esto permite desacoplar el Service de la implementación concreta de almacenamiento.
func NewService(repo BoardRepository, logger *zap.SugaredLogger) *Service {
	return &Service{
		Repo:   repo,
		Logger: logger,
	}
}

func (s *Service) Create(ctx context.Context, cBR CreateBoardRequest) (*CreateBoardResponse, error) {
	// Recibimos el payload del handler
	board, err := MapCreateBoardRequestToBoard(cBR)
	if err != nil {
		return nil, err
	}
	if err := s.Repo.Create(ctx, board); err != nil {
		return nil, err
	}

	resp := &CreateBoardResponse{
		ID:        board.ID,
		CreatedAt: board.CreatedAt,
	}

	return resp, nil
}

func (s *Service) Get(ctx context.Context) ([]BoardResponse, error) {
	boards, err := s.Repo.Get(ctx)
	if err != nil {
		return nil, err
	}

	resp := MapToBoardsResponse(boards)

	return resp, nil
}

func (s *Service) GetBoardsForUser(ctx context.Context, userID int64) ([]BoardForUsersResponse, error) {
	boards, err := s.Repo.GetBoardsForUser(ctx, userID)
	if err != nil {
		return nil, err
	}

	resp := MapToBoardsForUsersResponse(boards)

	return resp, nil
}

func (s *Service) GetById(ctx context.Context, id int64) (*BoardResponse, error) {
	board, err := s.Repo.GetById(ctx, id)
	if err != nil {
		return nil, err
	}

	resp := MapToBoardResponse(board)

	return resp, nil
}

func (s *Service) Update(ctx context.Context, id int64, uBR UpdateBoardRequest) (*UpdateBoardResponse, error) {
	// Obtenemos el usuario
	boardReq, err := s.Repo.GetById(ctx, id)
	if err != nil {
		return nil, err
	}
	// Verificamos que exista
	if boardReq == nil {
		return nil, utils.ErrNotFound
	}
	// Validamos y Pasamos a una funcion todos los datos a validar
	board, err := MapUpdateBoardRequestToBoard(uBR, boardReq)
	if err != nil {
		return nil, err
	}
	// Mandamos el user al almacenamiento
	if err := s.Repo.Update(ctx, board); err != nil {
		return nil, err
	}
	// Transformamos la respuesta al Response deseado
	resp := &UpdateBoardResponse{
		ID:        board.ID,
		UpdatedAt: board.UpdatedAt,
	}

	return resp, nil
}

func (s *Service) Delete(ctx context.Context, id int64) error {
	return s.Repo.Delete(ctx, id)
}
