package main

import (
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"time"

	"github.com/kacperhemperek/twitter-v2/api"
	"github.com/kacperhemperek/twitter-v2/router"
)

func init() {
	api.LoadEnv()
}

func main() {
	db := api.NewDB()
	r := router.New(db)
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
