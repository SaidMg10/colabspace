package user

import (
	"context"
	"errors"

	"github.com/SaidMg10/colabspace/internal/utils"
	"go.uber.org/zap"
)

// Service representa la capa de lógica de negocio para una entidad.
// Depende de una interfaz de repositorio que define las operaciones de persistencia necesarias.
type Service struct {
	// Se inyecta la interfaz del repositorio (UserRepository),
	// permitiendo que el Service use cualquier implementación concreta (por ejemplo, UserStore).
	Repo   UserRepository
	Logger *zap.SugaredLogger
}

// NewService es el constructor del Service.
// Recibe una implementación del repositorio y retorna una instancia del Service.
// Esto permite desacoplar el Service de la implementación concreta de almacenamiento.
func NewService(repo UserRepository, logger *zap.SugaredLogger) *Service {
	return &Service{
		Repo:   repo,
		Logger: logger,
	}
}

// Create crea un nuevo usuario a partir del DTO CreateUser.
// Recibe un contexto y los datos para crear el usuario.
// Transforma la contraseña, construye el modelo User,
// llama al repositorio para persistirlo y devuelve una respuesta o error.
func (s *Service) Create(ctx context.Context, createUserRequest CreateUserRequest) (*CreateUserResponse, error) {
	user, err := MapCreateUserRequestToUser(createUserRequest)
	if err != nil {
		return nil, err
	}
	if err := s.Repo.Create(ctx, user); err != nil {
		switch {
		case errors.Is(err, ErrDuplicateEmail):
			// manejar el error
			return nil, ErrDuplicateEmail
		case errors.Is(err, ErrDuplicateUsername):
			// manejar el error
			return nil, ErrDuplicateUsername
		default:
			return nil, err
		}
	}

	resp := &CreateUserResponse{
		ID:        user.ID,
		CreatedAt: user.CreatedAt,
	}

	return resp, nil
}

func (s *Service) Get(ctx context.Context) ([]UserSummaryResponse, error) {
	users, err := s.Repo.Get(ctx)
	if err != nil {
		return nil, err
	}

	resp := MapToUsersSummaryResponse(users)

	return resp, nil
}

func (s *Service) GetById(ctx context.Context, id int64) (*UserSummaryResponse, error) {
	user, err := s.Repo.GetById(ctx, id)
	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, utils.ErrNotFound
	}

	resp := MapToUserSummaryResponse(user)
	return resp, nil
}

func (s *Service) Update(ctx context.Context, id int64, updateUserRequest UpdateUserRequest) (*UpdateUserResponse, error) {
	// Obtenemos el usuario
	userReq, err := s.Repo.GetById(ctx, id)
	if err != nil {
		return nil, err
	}
	// Verificamos que exista
	if userReq == nil {
		return nil, utils.ErrNotFound
	}
	// Validamos y Pasamos a una funcion todos los datos a validar
	user, err := MapUpdateUserRequestToUser(updateUserRequest, userReq)
	if err != nil {
		return nil, err
	}
	// Mandamos el user al almacenamiento
	if err := s.Repo.Update(ctx, user); err != nil {
		return nil, err
	}
	// Transformamos la respuesta al Response deseado
	resp := &UpdateUserResponse{
		ID:        user.ID,
		UpdatedAt: user.UpdatedAt,
	}

	return resp, nil
}

func (s *Service) Delete(ctx context.Context, id int64) error {
	return s.Repo.Delete(ctx, id)
}

// Funciones independientes
func (s *Service) GetUserById(ctx context.Context, id int64) (*User, error) {
	user, err := s.Repo.GetById(ctx, id)
	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, utils.ErrNotFound
	}

	return user, nil
}

func (s *Service) GetByUsernameOrEmail(ctx context.Context, value string) (*User, error) {
	user, err := s.Repo.GetByUsernameOrEmail(ctx, value)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, utils.ErrNotFound
	}
	return user, nil
}
