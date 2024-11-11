package services

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/kacperhemperek/twitter-v2/api"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

type UserService struct {
	db neo4j.DriverWithContext
}

var ErrUserNotFound = errors.New("user not found")
var ErrInvalidUserQueryResponse = errors.New("invalid user query response")

func (s UserService) GetByEmail(ctx context.Context, email string) (user any, err error) {
	defer func() {
		if err != nil {
			LogServiceError("user", "find user by email", err)
		}
	}()

	q := `MATCH (user:User { email: $email }) RETURN user LIMIT 1`
	result, err := neo4j.ExecuteQuery(
		ctx,
		s.db,
		q,
		map[string]any{"email": email},
		neo4j.EagerResultTransformer,
		neo4j.ExecuteQueryWithDatabase(api.ENV.DB_NAME),
	)

	if err != nil {
		return nil, err
	}

	if len(result.Records) == 0 {
		return nil, ErrUserNotFound
	}

	user, ok := result.Records[0].Get("user")

	if !ok {
		return nil, ErrInvalidUserQueryResponse
	}

	return user, nil
}

func (s UserService) CreateUser(ctx context.Context, email, name string) (user any, err error) {
	defer func() {
		if err != nil {
			LogServiceError("user", "create a user", err)
		}
	}()
	q := `CREATE (user:User { email: $email, name: $name }) RETURN user LIMIT 1`
	result, err := neo4j.ExecuteQuery(
		ctx,
		s.db,
		q,
		map[string]any{"email": email, "name": name},
		neo4j.EagerResultTransformer,
		neo4j.ExecuteQueryWithDatabase(api.ENV.DB_NAME),
	)

	if err != nil {
		return nil, err
	}

	if len(result.Records) == 0 {
		err = errors.New("user not retrieved from create")
		return nil, err
	}

	user, ok := result.Records[0].Get("user")

	if !ok {
		return nil, ErrInvalidUserQueryResponse
	}

	return user, nil
}

func NewUserService(db neo4j.DriverWithContext) *UserService {
	return &UserService{
		db: db,
	}
}

func LogServiceError(service, method string, err error) {
	m := fmt.Sprintf("failed %s method", method)
	s := fmt.Sprintf("%s service", service)
	slog.Error(s, "message", m, "error", err)
}
