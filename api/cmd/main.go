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
)

func init() {
	api.LoadEnv()
	api.SetupLogger()
}

func main() {
	auth.Setup()
	db := api.NewDB()

	userService := services.NewUserService(db)
	sessionService := auth.NewSessionService(db)

	r := router.New(*userService, *sessionService)
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
