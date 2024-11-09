package api

import (
	"log/slog"
	"os"

	"github.com/joho/godotenv"
)

type env struct {
	FRONTEND_URL         string
	DB_PASSWORD          string
	DB_USERNAME          string
	DB_URL               string
	GOOGLE_CLIENT_ID     string
	GOOGLE_CLIENT_SECRET string
}

// NOTE: this will have all variables defined in the env file
var ENV *env = nil

func LoadEnv() {
	slog.Info("env", "message", "loading ENV into memory")
	err := godotenv.Load()

	if err != nil {
		slog.Error("env", "message", "could not load env from .env file, env will default to system env", "error", err)
	}

	frontendURL := loadVar("FRONTEND_URL", "http://localhost:3001")
	dbPassword := loadVar("DB_PASSWORD", "secret")
	dbUsername := loadVar("DB_USERNAME", "neo4j")
	dbUrl := loadVar("DB_URL", "bolt://neo4j:7687")
	googleClientID := loadVar("GOOGLE_CLIENT_ID", "google_client_id")
	googleClientSecret := loadVar("GOOGLE_CLIENT_SECRET", "google_secret")

	if ENV != nil {
		slog.Error("env", "message", "ENV already loaded")
		panic(1)
	}

	ENV = &env{
		FRONTEND_URL:         frontendURL,
		DB_PASSWORD:          dbPassword,
		DB_USERNAME:          dbUsername,
		DB_URL:               dbUrl,
		GOOGLE_CLIENT_ID:     googleClientID,
		GOOGLE_CLIENT_SECRET: googleClientSecret,
	}
}

func loadVar(name, def string) string {
	val := os.Getenv(name)

	if len(val) == 0 {
		return def
	}

	return val
}
