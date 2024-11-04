package api

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

func ApplyCors(r *mux.Router) http.Handler {
	opts := cors.Options{
		AllowedOrigins:   []string{ENV.FRONTEND_URL},
		AllowCredentials: true,
	}
	return cors.New(opts).Handler(r)
}
