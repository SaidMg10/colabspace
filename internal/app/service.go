package app

import (
	"github.com/SaidMg10/colabspace/internal/app/auth"
	"github.com/SaidMg10/colabspace/internal/app/user"
	"github.com/SaidMg10/colabspace/internal/config"
	"github.com/SaidMg10/colabspace/internal/storage"
	"go.uber.org/zap"
)

// Services es la estructura principal que agrupa todos los servicios de la aplicación.
// Cada campo representa un servicio asociado a una entidad del dominio.
// Ej: User user.Service, Post post.Service, etc.
type Services struct {
	// Servicio relacionado a la entidad User.
	// Su tipo corresponde al struct definido en user/services.go.
	User user.Service
	Auth auth.Service
}

// NewServices es el constructor del contenedor de servicios de la app.
// Recibe el Storage que agrupa todos los repositorios de entidades,
// e inicializa los servicios inyectando los repositorios correspondientes.
func NewServices(storage storage.Storage, logger *zap.SugaredLogger, authCfg config.TokenConfig) *Services {
	authenticator := auth.NewJWTAuthenticatorFromConfig(authCfg)
	return &Services{
		User: *user.NewService(storage.Users, logger),
		Auth: *auth.NewService(
			storage.Users,
			logger,
			authenticator,
			authCfg.Exp,
			authCfg.Iss,
		),
	}
}
