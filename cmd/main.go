package main

import (
	"github.com/SaidMg10/colabspace/internal/app"
	"github.com/SaidMg10/colabspace/internal/config"
	"github.com/SaidMg10/colabspace/internal/database"
	"go.uber.org/zap"
)

func main() {
	cfg := config.LoadConfig()

	logger := zap.Must(zap.NewProduction()).Sugar()
	defer logger.Sync()

	db, err := database.New(
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

	app := &app.Application{
		Config: &cfg,
		Logger: logger,
	}

	if err := app.Run(); err != nil {
		logger.Fatal(err)
	}
}
