package api

import (
	"log/slog"
	"os"

	"github.com/joho/godotenv"
)

type env struct {
	FRONTEND_URL string
	DB_PASSWORD  string
	DB_USERNAME  string
}

// NOTE: this will have all variables defined in the env file
var ENV *env = nil

func LoadEnv() {
	slog.Info("env", "message", "loading ENV into memory")
	err := godotenv.Load()

	if err != nil {
		slog.Error("env", "load err", err)
		panic(1)
	}

	frontendURL := loadVar("FRONTEND_URL", "http://localhost:3001")
	dbPassword := loadVar("DB_PASSWORD", "secret")
	dbUsername := loadVar("DB_USERNAME", "neo4j")

	if ENV != nil {
		slog.Error("env", "message", "ENV already loaded")
		return
	}

	ENV = &env{
		FRONTEND_URL: frontendURL,
		DB_PASSWORD:  dbPassword,
		DB_USERNAME:  dbUsername,
	}
}

func loadVar(name, def string) string {
	val := os.Getenv(name)

	if len(val) == 0 {
		return def
	}

	return val
}
