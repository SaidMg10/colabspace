package boardusers

import (
	"context"

	"github.com/SaidMg10/colabspace/internal/app/user"
	"go.uber.org/zap"
)

type Service struct {
	Repo     BoardUsersRepository
	UserRepo user.UserRepository
	Logger   *zap.SugaredLogger
}

// NewService es el constructor del Service.
// Recibe una implementación del repositorio y retorna una instancia del Service.
// Esto permite desacoplar el Service de la implementación concreta de almacenamiento.
func NewService(repo BoardUsersRepository, userRepo user.UserRepository, logger *zap.SugaredLogger) *Service {
	return &Service{
		Repo:     repo,
		UserRepo: userRepo,
		Logger:   logger,
	}
}

func (s *Service) Add(ctx context.Context, boardUsersRequest *CreateBoardUsersRequest) error {
	user, err := s.UserRepo.GetByUsernameOrEmail(ctx, boardUsersRequest.Value)
	if err != nil {
		return err
	}
	boardUsers, err := MapCreateBoardUsersRequestToBoardUsers(boardUsersRequest.BoardID, user.ID)
	if err != nil {
		return err
	}
	if err := s.Repo.Add(ctx, boardUsers); err != nil {
		return err
	}
	return nil
}

func (s *Service) GetMembers(ctx context.Context, boardID int64) ([]user.UserSummaryResponse, error) {
	users, err := s.Repo.GetMembers(ctx, boardID)
	if err != nil {
		return nil, err
	}

	resp := user.MapToUsersSummaryResponse(users)
	return resp, nil
}

func (s *Service) Leave(ctx context.Context, boardID, userID int64) error {
	if err := s.Repo.Leave(ctx, boardID, userID); err != nil {
		return err
	}
	return nil
}

func (s *Service) Delete(ctx context.Context, boardID, userID int64) error {
	if err := s.Repo.Delete(ctx, boardID, userID); err != nil {
		return err
	}
	return nil
}

func (s *Service) IsMember(ctx context.Context, boardID, userID int64) (bool, error) {
	isMember, err := s.Repo.IsMember(ctx, boardID, userID)
	if err != nil {
		return false, err
	}
	return isMember, nil
}
