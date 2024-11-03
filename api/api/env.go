package api

import (
	"log/slog"
	"os"

	"github.com/joho/godotenv"
)

type env struct {
	FRONTEND_URL string
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

	if ENV != nil {
		slog.Error("env", "message", "ENV already loaded")
		return
	}

	ENV = &env{
		FRONTEND_URL: frontendURL,
	}
}

func loadVar(name, def string) string {
	val := os.Getenv(name)

	if len(val) == 0 {
		return def
	}

	return val
}
