package http

import (
	"github.com/SaidMg10/colabspace/internal/app"
	"github.com/SaidMg10/colabspace/internal/app/auth"
	"github.com/SaidMg10/colabspace/internal/transport/http/handlers"
	"github.com/SaidMg10/colabspace/internal/transport/http/middleware"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// NewRouter configura y devuelve una nueva instancia del router de Gin.
// Se espera que reciba una instancia de `app.Services`, que contiene los servicios
// que necesitan los handlers para funcionar.
func NewRouter(s *app.Services, logger *zap.SugaredLogger, authenticator auth.Authenticator) *gin.Engine {
	// Creamos una nueva instancia del router usando Gin.
	// `gin.New()` crea un router limpio (sin middlewares por defecto).
	router := gin.New()
	// Aplicamos middlewares globales al router.
	// `gin.Recovery()` evita que un panic haga caer la app; devuelve 500 si algo explota.
	router.Use(gin.Recovery()) // Recupera de panics y evita que caiga el servidor

	middleware := middleware.NewMiddleware(authenticator, s)

	// La func .Group() es para agregar un prefix y para agrupar rutas
	// Ej: test := router.Group("testing") -> testv1 := test.Group("testv1")
	// Creamos un grupo de rutas con el prefijo /api
	api := router.Group("api")
	v1 := api.Group("v1")

	// Registra rutas de los handlers
	healthHandler := handlers.NewHealthHandler()
	// Registramos una ruta GET que responde en /api/v1/health

	v1.GET("/health", healthHandler.Check)
	// Registro de Auth
	authHandler := handlers.NewAuthHandler(&s.Auth, logger)
	authHandler.RegisterRoutes(v1)
	// Aquí podrías registrar más handlers, como por ejemplo usuarios:
	// Registro de User
	v1.Use(middleware.AuthTokenMiddleware())
	userHandler := handlers.NewUserHandler(&s.User, logger)
	userHandler.RegisterRoutes(v1)
	// Registro de Board
	boardHandler := handlers.NewBoardHandler(&s.Board, logger, middleware)
	boardHandler.RegisterRoutes(v1)
	// Registro de Añadir Usuarios
	boardUsersHandler := handlers.NewBoardUsersHandler(&s.BoardUsers, logger, middleware)
	boardUsersHandler.RegisterRoutes(v1)
	// Registro de Stages
	stagesHandler := handlers.NewStageHandler(&s.Stage, logger, middleware)
	stagesHandler.RegisterRoutes(v1)

	return router
}
