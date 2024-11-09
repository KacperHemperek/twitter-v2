package router

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/kacperhemperek/twitter-v2/api"
	"github.com/kacperhemperek/twitter-v2/auth"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

func New(db neo4j.DriverWithContext) *mux.Router {
	r := mux.NewRouter()

	r.HandleFunc("/api/healthcheck", api.Handle(func(w http.ResponseWriter, r *http.Request) error {
		return api.JSON(w, map[string]any{"message": "OK"}, http.StatusOK)
	})).Methods(http.MethodGet)

	r.HandleFunc("/api/auth/{provider}/login", api.Handle(auth.LoginHandler())).Methods(http.MethodGet)
	r.HandleFunc("/api/auth/{provider}/callback", api.Handle(auth.AuthCallbackHanlder())).Methods(http.MethodGet)
	r.HandleFunc("/api/auth/{provider}/login", api.Handle(auth.LogoutHandler())).Methods(http.MethodGet)

	return r
}
