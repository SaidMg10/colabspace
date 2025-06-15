package server

import (
	"log"

	"github.com/SaidMg10/colabspace/internal/env"
	"github.com/SaidMg10/colabspace/internal/store"
	"github.com/joho/godotenv"
	"go.uber.org/zap"
)

type Application struct {
	// TODO: Add logger, store, etc.
	Config Config
	Logger *zap.SugaredLogger
	Store  store.Storage
}

type Config struct {
	Addr        string
	Db          DbConfig
	Env         string
	ApiURL      string
	FrontendURL string
}

type DbConfig struct {
	Addr         string
	MaxOpenConns int
	MaxIdleConns int
	MaxIdleTime  string
}

// LoadConfig carga la Configuración desde variables de entorno o usa valores por defecto
func LoadConfig() Config {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables directly")
	}

	cfg := Config{
		Addr:        env.GetString("ADDR", ":8080"),
		ApiURL:      env.GetString("EXTERNAL_URL", "localhost:8080"),
		FrontendURL: env.GetString("FRONTEND_URL", "http://localhost:5173"),
		Db: DbConfig{
			Addr:         env.GetString("DB_ADDR", "postgres://admin:adminpassword@localhost/colabs?sslmode=disable"),
			MaxOpenConns: env.GetInt("DB_MAX_OPEN_CONNS", 30),
			MaxIdleConns: env.GetInt("DB_MAX_IDLE_CONNS", 30),
			MaxIdleTime:  env.GetString("DB_MAX_IDLE_TIME", "15m"),
		},
		Env: env.GetString("ENV", "development"),
	}
	return cfg
}
