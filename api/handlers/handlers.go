package handlers

import (
	"log/slog"
	"net/http"

	"github.com/kacperhemperek/twitter-v2/api"
	"github.com/kacperhemperek/twitter-v2/services"
)

func CreateTweetHandler(tweetService services.TweetService) api.HandlerFunc {
	type request struct {
		Body string `json:"body" validate:"required,max=180,min=1"`
	}
	type response struct {
		ID string `json:"id"`
	}
	return func(w http.ResponseWriter, r *api.Request) error {
		body := &request{}
		if err := r.ValidateBody(body); err != nil {
			return err
		}
		user, _ := r.User()
		slog.Debug("creating tweet", "user", user, "tweetBody", body.Body)
		tweet, err := tweetService.Create(r.Context(), user.ID, body.Body)
		if err != nil {
			return err
		}
		return api.JSON(w, &response{ID: tweet.ID}, http.StatusCreated)
	}
}
