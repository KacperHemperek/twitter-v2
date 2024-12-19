package store

import (
	"errors"
	"log/slog"

	"github.com/kacperhemperek/twitter-v2/api"
	"github.com/kacperhemperek/twitter-v2/lib/dbmap"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j/dbtype"
)

var (
	ErrInvalidResponseType = errors.New("unsupported response type")
)

type Store = neo4j.DriverWithContext

// New creates a new connection to the database and returns it, panics if connection could not be established
func New() Store {
	slog.Info("db", "message", "requesting connection to the database")
	driver, err := neo4j.NewDriverWithContext(api.ENV.DB_URL, neo4j.BasicAuth(api.ENV.DB_USERNAME, api.ENV.DB_PASSWORD, ""))

	if err != nil {
		slog.Error("db", "message", "could not establish connection to the database", "err", err)
		panic(1)
	}

	slog.Info("db", "message", "connected to the database")
	return driver
}

// WithApiStore returns a configuration option for executing queries with the default api database
func WithAPIStore() neo4j.ExecuteQueryConfigurationOption {
	return neo4j.ExecuteQueryWithDatabase(api.ENV.DB_NAME)
}

// Read reads the result of a query from the database response
// and decodes it into a struct provided in the val argument
// raw should be of type dbtype.Node or dbtype.Relationship but is used here as any to
// avoid problems with types in services
//
// val should be a pointer to the struct that the result should be decoded into
func Read(raw any, val any) error {
	switch userNode := raw.(type) {
	case dbtype.Node:
		err := dbmap.Decode(userNode.GetProperties(), val)
		if err != nil {
			return err
		}
		return nil
	case dbtype.Relationship:
		err := dbmap.Decode(userNode.GetProperties(), val)
		if err != nil {
			return err
		}
		return nil
	default:
		return ErrInvalidResponseType
	}
}
