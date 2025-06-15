package main

import (
	"github.com/SaidMg10/colabspace/internal/api/router"
	db "github.com/SaidMg10/colabspace/internal/database"
	"github.com/SaidMg10/colabspace/internal/server"
	"github.com/SaidMg10/colabspace/internal/store"
	"go.uber.org/zap"
)

func main() {
	cfg := server.LoadConfig()
	logger := zap.Must(zap.NewProduction()).Sugar()
	defer logger.Sync()

	db, err := db.New(cfg.Db.Addr, cfg.Db.MaxOpenConns, cfg.Db.MaxIdleConns, cfg.Db.MaxIdleTime)
	if err != nil {
		logger.Fatal(err)
	}
	defer db.Close()
	logger.Info("database connection pool established")

	store := store.NewStorage(db)

	app := &server.Application{
		Config: cfg,
		Logger: logger,
		Store:  store,
	}

	mux := router.Mount(app)

	logger.Fatal(app.Run(mux))
}
