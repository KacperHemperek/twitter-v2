package router

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/kacperhemperek/twitter-v2/api"
)

func New() *mux.Router {
	r := mux.NewRouter()

	r.HandleFunc("/healthcheck", api.Handle(func(w http.ResponseWriter, r *http.Request) error {
		return api.JSON(w, map[string]any{"message": "OK"}, http.StatusOK)
	})).Methods(http.MethodGet)

	return r
}
