package main

import (
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"time"

	"github.com/kacperhemperek/twitter-v2/api"
	"github.com/kacperhemperek/twitter-v2/auth"
	"github.com/kacperhemperek/twitter-v2/router"
	"github.com/kacperhemperek/twitter-v2/services"
	"github.com/kacperhemperek/twitter-v2/store"
)

func init() {
	api.LoadEnv()
	api.SetupLogger()
}

func main() {
	slog.Info("main", "message", "starting application", "environment", api.ENV.ENVIRONMENT)
	auth.Setup()
	db := store.New()

	userService := services.NewUserService(db)
	sessionService := auth.NewSessionService(db)
	tweetService := services.NewTweetService(db)

	h := api.NewAPIHandler()

	r := router.New(h, *userService, *sessionService, *tweetService)
	handler := api.ApplyCors(r)

	p := 1337
	server := &http.Server{
		Handler:      handler,
		Addr:         fmt.Sprintf("0.0.0.0:%d", p),
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	slog.Info("main", "message", fmt.Sprintf("server listening on port :%d", p))
	log.Fatalln(server.ListenAndServe())
}
