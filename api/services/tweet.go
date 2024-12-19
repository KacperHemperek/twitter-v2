package services

import (
	"context"

	"github.com/kacperhemperek/twitter-v2/api"
	"github.com/kacperhemperek/twitter-v2/models"
	"github.com/kacperhemperek/twitter-v2/store"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

type TweetService struct {
	db store.Store
}

// Create creates a new tweet and returns it if creation was successful, otherwise returns an error
func (s *TweetService) Create(ctx context.Context, authorID, body string) (ct *models.TweetModel, err error) {
	defer func() {
		if err != nil {
			api.LogServiceError("tweet", "create tweet", err)
		}
	}()

	q := `
	MATCH (user:User { id: $authorId })
	CREATE (user)-[:TWEETED {createdAt: datetime()}]->(tweet:Tweet { id: randomUUID(), body: $body, createdAt: datetime(), authorId: $authorId })
	RETURN tweet
	`
	result, err := neo4j.ExecuteQuery(
		ctx,
		s.db,
		q,
		map[string]any{"authorId": authorID, "body": body},
		neo4j.EagerResultTransformer,
		store.WithAPIStore(),
	)

	if err != nil {
		return nil, err
	}

	if len(result.Records) == 0 {
		return nil, ErrResourceNotFound
	}

	rawUser, ok := result.Records[0].Get("tweet")

	if !ok {
		return nil, ErrInvalidQueryResponse
	}

	t := &models.TweetModel{}
	if err := store.Read(rawUser, t); err != nil {
		return nil, err
	}
	return t, nil
}

// GetByID returns a tweet by its ID if it exists, otherwise returns an error
func (s *TweetService) GetByID(ctx context.Context, id string) (t *models.TweetModel, err error) {
	return nil, nil
}

func NewTweetService(db neo4j.DriverWithContext) *TweetService {
	return &TweetService{
		db: db,
	}
}
