package config

import (
	"log"
	"time"

	"github.com/SaidMg10/colabspace/internal/env"
	"github.com/joho/godotenv"
)

type Config struct {
	Server   ServerConfig
	Database DbConfig
	App      AppConfig
}

type ServerConfig struct {
	Host            string        // IP o hostname donde correrá el servidor, ej: "0.0.0.0" o "localhost"
	Port            int           // Puerto en el que escuchará el servidor, ej: 8080
	ReadTimeout     time.Duration // Timeout para leer la request HTTP
	WriteTimeout    time.Duration // Timeout para escribir la response HTTP
	IdleTimeout     time.Duration // Timeout para conexiones inactivas
	ShutdownTimeout time.Duration // Timeout para un apagado ordenado del servidor
}

type DbConfig struct {
	DSN          string
	MaxOpenConns int
	MaxIdleConns int
	MaxIdleTime  time.Duration
}

type AppConfig struct {
	Env         string
	ApiURL      string
	FrontendURL string
	Auth        AuthConfig
}

type AuthConfig struct {
	Basic BasicConfig
	Token TokenConfig
}

type TokenConfig struct {
	Secret string
	Exp    time.Duration
	Iss    string
}

type BasicConfig struct {
	User string
	Pass string
}

func LoadConfig() Config {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables directly.")
	}
	return Config{
		Server: ServerConfig{
			Host:            env.GetString("SERVER_HOST", "0.0.0.0"),
			Port:            env.GetInt("SERVER_PORT", 8080),
			ReadTimeout:     env.GetDuration("SERVER_READ_TIMEOUT", 5*time.Second),
			WriteTimeout:    env.GetDuration("SERVER_WRITE_TIMEOUT", 10*time.Second),
			IdleTimeout:     env.GetDuration("SERVER_IDLE_TIMEOUT", 120*time.Second),
			ShutdownTimeout: env.GetDuration("SERVER_SHUTDOWN_TIMEOUT", 10*time.Second),
		},
		Database: DbConfig{
			DSN:          env.GetString("DB_DSN", ""),
			MaxOpenConns: env.GetInt("DB_MAX_OPEN_CONNS", 10),
			MaxIdleConns: env.GetInt("DB_MAX_IDLE_CONNS", 5),
			MaxIdleTime:  env.GetDuration("DB_MAX_IDLE_TIME", 15*time.Minute),
		},
		App: AppConfig{
			Env:         env.GetString("ENV", "development"),
			ApiURL:      env.GetString("API_URL", ""),
			FrontendURL: env.GetString("FRONTEND_URL", ""),
			Auth: AuthConfig{
				Basic: BasicConfig{
					User: env.GetString("AUTH_BASIC_USER", ""),
					Pass: env.GetString("AUTH_BASIC_PASS", ""),
				},
				Token: TokenConfig{
					Secret: env.GetString("AUTH_TOKEN_SECRET", "Pacoco"),
					Exp:    env.GetDuration("AUTH_TOKEN_EXP", time.Hour*24),
					Iss:    env.GetString("AUTH_TOKEN_ISS", "myapp"),
				},
			},
		},
	}
}
