package auth

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"time"

	"github.com/go-viper/mapstructure/v2"
	"github.com/kacperhemperek/twitter-v2/api"
	"github.com/kacperhemperek/twitter-v2/models"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j/dbtype"
)

var (
	ErrSessionNotCreate          = errors.New("could not create session")
	ErrInvalidSessionQueryResult = errors.New("invalid session query response")
	ErrSessionNotFound           = errors.New("session not found")
)

type SessionService struct {
	store neo4j.DriverWithContext
}

func (s SessionService) GetSession(ctx context.Context, ID string) (sess *models.SessionModel, err error) {
	defer func() {
		if err != nil {
			api.LogServiceError("session", "get session", err)
		}
	}()

	q := `
	MATCH (session:Session { id: $sessionID })
	RETURN session
	LIMIT 1;
	`
	result, err := neo4j.ExecuteQuery(
		ctx,
		s.store,
		q,
		map[string]any{"sessionID": ID},
		neo4j.EagerResultTransformer,
		neo4j.ExecuteQueryWithDatabase(api.ENV.DB_NAME),
	)

	if err != nil {
		return nil, err
	}

	if len(result.Records) == 0 {
		return nil, ErrSessionNotFound
	}

	rawSession, ok := result.Records[0].Get("session")
	if !ok {
		return nil, ErrInvalidSessionQueryResult
	}

	switch sessionNode := rawSession.(type) {
	case dbtype.Node:
		slog.Info("session service", "rawSession", sessionNode.GetProperties())
		err = mapstructure.Decode(sessionNode.GetProperties(), &sess)
		if err != nil {
			return nil, ErrInvalidSessionQueryResult
		}
		return sess, nil
	default:
		return nil, ErrInvalidSessionQueryResult
	}
}

func (s SessionService) DeleteSession(ctx context.Context, sessID string) (err error) {
	defer func() {
		if err != nil {
			api.LogServiceError("session", "delete session", err)
		}
	}()

	q := `
	MATCH (session:Session { id: $ID }) DETACH DELETE session;
	`
	_, err = neo4j.ExecuteQuery(
		ctx,
		s.store,
		q,
		map[string]any{"ID": sessID},
		neo4j.EagerResultTransformer,
		neo4j.ExecuteQueryWithDatabase(api.ENV.DB_NAME),
	)

	return err
}

func (s SessionService) CreateSession(ctx context.Context, userID string) (sess *models.SessionModel, err error) {
	defer func() {
		if err != nil {
			api.LogServiceError("session", "create session", err)
		}
	}()

	q := `
	MATCH (user:User { id: $userID })
	CREATE (session:Session { id: randomUUID(), expiration: datetime($exp), userId: $userID })<-[:CREATED]-(user)
	RETURN session LIMIT 1
	`
	exp := time.Now().Add(time.Hour * 24 * 30).Format(time.RFC3339)
	result, err := neo4j.ExecuteQuery(
		ctx,
		s.store,
		q,
		map[string]any{"userID": userID, "exp": exp},
		neo4j.EagerResultTransformer,
		neo4j.ExecuteQueryWithDatabase(api.ENV.DB_NAME),
	)

	if err != nil {
		return nil, err
	}

	if len(result.Records) == 0 {
		return nil, ErrInvalidSessionQueryResult
	}

	rawSession, ok := result.Records[0].Get("session")
	if !ok {
		return nil, ErrInvalidSessionQueryResult
	}
	switch sessionNode := rawSession.(type) {
	case dbtype.Node:
		err = mapstructure.Decode(sessionNode.GetProperties(), &sess)
		if err != nil {
			return nil, ErrInvalidSessionQueryResult
		}
		return sess, nil
	default:
		return nil, ErrInvalidSessionQueryResult
	}
}

func NewSessionService(store neo4j.DriverWithContext) *SessionService {
	return &SessionService{
		store: store,
	}
}

func SetSessionCookie(w http.ResponseWriter, sessID string) {
	http.SetCookie(w, newSessionCookie(sessID, 60*60*24*30))
}

func ClearSessionCookie(w http.ResponseWriter) {
	http.SetCookie(w, newSessionCookie("", -1))
}

func newSessionCookie(sessID string, maxAge int) *http.Cookie {
	var ss http.SameSite
	if api.ENV.IsProd() {
		ss = http.SameSiteStrictMode
	} else {
		ss = http.SameSiteDefaultMode
	}
	return &http.Cookie{
		Name:     "sessionID",
		Value:    sessID,
		Path:     "/",
		MaxAge:   maxAge,
		Secure:   api.ENV.IsProd(),
		HttpOnly: true,
		SameSite: ss,
	}
}
