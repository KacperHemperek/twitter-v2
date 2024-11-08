package api

import (
	"log/slog"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

func NewDB() neo4j.DriverWithContext {
	slog.Info("db", "message", "requesting connection to the database")
	driver, err := neo4j.NewDriverWithContext(ENV.DB_URL, neo4j.BasicAuth(ENV.DB_USERNAME, ENV.DB_PASSWORD, ""))

	if err != nil {
		slog.Error("db", "message", "could not establish connection to the database", "err", err)
		panic(1)
	}

	slog.Info("db", "message", "connected to the database")
	return driver
}
