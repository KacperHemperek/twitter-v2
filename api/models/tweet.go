package models

import (
	"time"

	"github.com/kacperhemperek/twitter-v2/lib/dbmap"
)

type TweetModel struct {
	ID        string        `json:"id" dbmap:"id"`
	Body      string        `json:"body" dbmap:"body"`
	CreatedAt time.Time     `json:"createdAt" dbmap:"createdAt"`
	DeletedAt dbmap.NilTime `json:"deletedAt" dbmap:"deletedAt"`
	AuthorID  string        `json:"authorId" dbmap:"authorId"`
}
