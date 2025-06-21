package main

import (
	"github.com/SaidMg10/colabspace/internal/app"
	"github.com/SaidMg10/colabspace/internal/app/auth"
	"github.com/SaidMg10/colabspace/internal/config"
	"github.com/SaidMg10/colabspace/internal/database"
	"github.com/SaidMg10/colabspace/internal/storage"
	"github.com/SaidMg10/colabspace/internal/transport/http"
	"go.uber.org/zap"
)

func main() {
	cfg := config.LoadConfig() // Cargamos Config

	logger := zap.Must(zap.NewProduction()).Sugar() // Cargamos Logger
	defer logger.Sync()

	db, err := database.New( // Cargamos Database
		cfg.Database.DSN,
		cfg.Database.MaxOpenConns,
		cfg.Database.MaxIdleConns,
		cfg.Database.MaxIdleTime,
	)
	if err != nil {
		logger.Fatal(err)
	}
	defer db.Close()
	logger.Info("database connection pool established")

	auth := auth.NewJWTAuthenticatorFromConfig(cfg.App.Auth.Token)

	storage := storage.NewStorage(db)                               // Cargamos Storage
	service := app.NewServices(storage, logger, cfg.App.Auth.Token) // Cargamos Service y se carga el store, logger y la config token en el service
	router := http.NewRouter(service, logger, auth)                 // Cargamos el Router y se carga el service y el logger en el router

	app := &app.Application{
		Config:  &cfg,
		Logger:  logger,
		Store:   storage,
		Service: *service,
		Router:  router,
	}

	if err := app.Run(); err != nil {
		logger.Fatal(err)
	}
}
