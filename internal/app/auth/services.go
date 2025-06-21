package auth

import (
	"context"
	"errors"
	"time"

	"github.com/SaidMg10/colabspace/internal/app/user"
	"github.com/golang-jwt/jwt/v5"
	"go.uber.org/zap"
)

// Creacion del struct del Service de auth
type Service struct {
	Repo          user.UserRepository // Repositorio del user
	Logger        *zap.SugaredLogger  // Logger
	Authenticator Authenticator       // LLamo al authenticator para pasar las interfaces
	TokenExpiry   time.Duration
	TokenIssuer   string
}

// Creacion de la funcion NewService que es el constructor del servicio
func NewService(
	repo user.UserRepository,
	logger *zap.SugaredLogger,
	authenticator Authenticator,
	expiry time.Duration,
	issuer string,
) *Service {
	return &Service{
		Repo:          repo,
		Logger:        logger,
		Authenticator: authenticator,
		TokenExpiry:   expiry,
		TokenIssuer:   issuer,
	}
}

// Crear función Login
func (s *Service) Login(ctx context.Context, login Login) (string, error) {
	// Se busca el usuario con el metodo GetByEmail (El usuario debe estar activo)
	user, err := s.Repo.GetByEmail(ctx, login.Email)
	if err != nil {
		return "", err // No se encontro el usuario o esta inactivo
	}
	// Ahora se compara la contraseña para poder continuar con el login
	err = user.Password.Compare(login.Password)
	if err != nil {
		return "", err // Contraseña incorrecta
	}
	// Generamos el claimsMap para los datos que recibira el token
	claims := jwt.MapClaims{
		"sub": user.ID,
		"rol": user.Role,
		"exp": time.Now().Add(s.TokenExpiry).Unix(), // Expiración (obligatoria si vas a validarla)
		"iat": time.Now().Unix(),                    // Fecha de emisión
		"iss": s.TokenIssuer,                        // Emisor
		"aud": s.TokenIssuer,                        // Audiencia
	}
	// Generamos el token pasando los cleims anteriormente mapeados
	token, err := s.Authenticator.GenerateToken(claims)
	if err != nil {
		return "", err
	}
	// Devolvemos el token
	return token, nil
}

func (s *Service) Register(ctx context.Context, createUserRequest user.CreateUserRequest) (*user.CreateUserResponse, error) {
	u, err := user.MapCreateUserRequestToUser(createUserRequest)
	if err != nil {
		return nil, err
	}
	if err := s.Repo.Create(ctx, u); err != nil {
		switch {
		case errors.Is(err, user.ErrDuplicateEmail):
			// manejar el error
			return nil, user.ErrDuplicateEmail
		case errors.Is(err, user.ErrDuplicateUsername):
			// manejar el error
			return nil, user.ErrDuplicateUsername
		default:
			return nil, err
		}
	}

	resp := &user.CreateUserResponse{
		ID:        u.ID,
		CreatedAt: u.CreatedAt,
	}

	return resp, nil
}
