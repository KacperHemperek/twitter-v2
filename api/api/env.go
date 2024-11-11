package api

import (
	"log/slog"
	"os"

	"github.com/joho/godotenv"
)

type env struct {
	FRONTEND_URL         string
	API_URL              string
	DB_PASSWORD          string
	DB_USERNAME          string
	DB_URL               string
	DB_NAME              string
	GOOGLE_CLIENT_ID     string
	GOOGLE_CLIENT_SECRET string
	SESSION_SECRET       string
	JWT_SECRET           string
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
	apiURL := loadVar("API_URL", "http://localhost:1337")
	dbPassword := loadVar("DB_PASSWORD", "secret")
	dbUsername := loadVar("DB_USERNAME", "neo4j")
	dbUrl := loadVar("DB_URL", "bolt://neo4j:7687")
	dbName := loadVar("DB_NAME", "neo4j")
	googleClientID := loadVar("GOOGLE_CLIENT_ID", "google_client_id")
	googleClientSecret := loadVar("GOOGLE_CLIENT_SECRET", "google_secret")
	sessionSecret := loadVar("SESSION_SECRET", "session_secret123")
	jwtSecret := loadVar("SESSION_SECRET", "jwt_secret123")

	if ENV != nil {
		slog.Error("env", "message", "ENV already loaded")
		panic(1)
	}

	ENV = &env{
		FRONTEND_URL:         frontendURL,
		API_URL:              apiURL,
		DB_PASSWORD:          dbPassword,
		DB_USERNAME:          dbUsername,
		DB_URL:               dbUrl,
		DB_NAME:              dbName,
		GOOGLE_CLIENT_ID:     googleClientID,
		GOOGLE_CLIENT_SECRET: googleClientSecret,
		SESSION_SECRET:       sessionSecret,
		JWT_SECRET:           jwtSecret,
	}
}

func loadVar(name, def string) string {
	val := os.Getenv(name)

	if len(val) == 0 {
		return def
	}

	return val
}
