package router

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/kacperhemperek/twitter-v2/api"
	"github.com/kacperhemperek/twitter-v2/auth"
	"github.com/kacperhemperek/twitter-v2/handlers"
	"github.com/kacperhemperek/twitter-v2/services"
)

func New(h *api.APIHandler, userService services.UserService, sessionService auth.SessionService) *mux.Router {
	r := mux.NewRouter()
	authMiddleware := auth.NewAuthMiddleware(userService, sessionService)

	r.HandleFunc("/api/healthcheck", h.Handle(func(w http.ResponseWriter, r *api.Request) error {
		return api.JSON(w, map[string]any{"message": "OK"}, http.StatusOK)
	})).Methods(http.MethodGet)


	r.HandleFunc("/api/auth/{provider}/login", h.Handle(auth.LoginHandler(userService, sessionService))).Methods(http.MethodGet)
	r.HandleFunc("/api/auth/{provider}/callback", h.Handle(auth.AuthCallbackHanlder(userService, sessionService))).Methods(http.MethodGet)
	r.HandleFunc("/api/auth/logout", h.Handle(authMiddleware(auth.LogoutHandler(sessionService)))).Methods(http.MethodGet)

	r.HandleFunc("/api/auth/me", h.Handle(authMiddleware(auth.GetMeHandler()))).Methods(http.MethodGet)

	return r
}
