package api

import (
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

var frontendURL = os.Getenv("FRONTEND_URL")

func ApplyCors(r *mux.Router) http.Handler {
	opts := cors.Options{
		AllowedOrigins:   []string{frontendURL},
		AllowCredentials: true,
	}
	return cors.New(opts).Handler(r)
}
