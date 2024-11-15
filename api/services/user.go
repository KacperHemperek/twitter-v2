package services

import (
	"context"
	"errors"

	"github.com/go-viper/mapstructure/v2"
	"github.com/kacperhemperek/twitter-v2/api"
	"github.com/kacperhemperek/twitter-v2/models"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j/dbtype"
)

type UserService struct {
	db neo4j.DriverWithContext
}

var ErrUserNotFound = errors.New("user not found")
var ErrInvalidUserQueryResponse = errors.New("invalid user query response")

func (s UserService) GetByEmail(ctx context.Context, email string) (user *models.UserModel, err error) {
	defer func() {
		if err != nil {
			api.LogServiceError("user", "find user by email", err)
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

	rawUser, ok := result.Records[0].Get("user")

	if !ok {
		return nil, ErrInvalidUserQueryResponse
	}

	switch userNode := rawUser.(type) {
	case dbtype.Node:
		err := mapstructure.Decode(userNode.GetProperties(), &user)
		if err != nil {
			return nil, ErrInvalidUserQueryResponse
		}
		return user, nil
	default:
		return nil, ErrInvalidUserQueryResponse
	}
}

func (s UserService) GetByID(ctx context.Context, ID string) (user *models.UserModel, err error) {
	defer func() {
		if err != nil {
			api.LogServiceError("user", "find user by email", err)
		}
	}()

	q := `MATCH (user:User { id: $ID }) RETURN user LIMIT 1`
	result, err := neo4j.ExecuteQuery(
		ctx,
		s.db,
		q,
		map[string]any{"ID": ID},
		neo4j.EagerResultTransformer,
		neo4j.ExecuteQueryWithDatabase(api.ENV.DB_NAME),
	)

	if err != nil {
		return nil, err
	}

	if len(result.Records) == 0 {
		return nil, ErrUserNotFound
	}

	rawUser, ok := result.Records[0].Get("user")

	if !ok {
		return nil, ErrInvalidUserQueryResponse
	}

	switch userNode := rawUser.(type) {
	case dbtype.Node:
		err := mapstructure.Decode(userNode.GetProperties(), &user)
		if err != nil {
			return nil, ErrInvalidUserQueryResponse
		}
		return user, nil
	default:
		return nil, ErrInvalidUserQueryResponse
	}
}

func (s UserService) CreateUser(ctx context.Context, email, name, avatar string) (user *models.UserModel, err error) {
	defer func() {
		if err != nil {
			api.LogServiceError("user", "create a user", err)
		}
	}()
	q := `
    CREATE (user:User {
        id: randomUUID(),
        email: $email,
        name: $name,
        image: $image,
        background: $background
    })
    RETURN user LIMIT 1
    `
	result, err := neo4j.ExecuteQuery(
		ctx,
		s.db,
		q,
		map[string]any{
			"email":      email,
			"name":       name,
			"image":      avatar,
			"background": "",
		},
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

	rawUser, ok := result.Records[0].Get("user")

	if !ok {
		return nil, ErrInvalidUserQueryResponse
	}

	switch userNode := rawUser.(type) {
	case dbtype.Node:
		err := mapstructure.Decode(userNode.GetProperties(), &user)
		if err != nil {
			return nil, ErrInvalidUserQueryResponse
		}
		return user, nil
	default:
		return nil, ErrInvalidUserQueryResponse
	}
}

func NewUserService(db neo4j.DriverWithContext) *UserService {
	return &UserService{
		db: db,
	}
}
