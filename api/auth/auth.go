package auth

import (
	"log/slog"
	"net/http"

	"github.com/kacperhemperek/twitter-v2/api"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/google"
)

func Setup() {
	slog.Info("auth", "message", "setting up oauth")
	goth.UseProviders(
		google.New(api.ENV.GOOGLE_CLIENT_ID, api.ENV.GOOGLE_CLIENT_SECRET, "http://localhost:1337/api/auth/google/callback"),
	)
	slog.Info("auth", "message", "oauth correctly setup")
}

// TODO: implement auth with google provider at least

func AuthCallbackHanlder() api.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		// TODO: i think should set something that tells frontend that user is
		// authenticated and that api can then read on request to get user info

		// TODO: create user in the database that can be then used if it does not exist yet
		user, err := gothic.CompleteUserAuth(w, r)
		if err != nil {
			return err
		}
		slog.Info("auth", "message", "user logged in successfully", "user", user)
		return nil
	}
}

func LogoutHandler() api.HandlerFunc {
	type response struct {
		Message string `json:"message"`
	}
	return func(w http.ResponseWriter, r *http.Request) error {
		// TODO: actually logout user (remove session/token or w/e)
		gothic.Logout(w, r)
		res := &response{
			Message: "user logged out successfully",
		}
		return api.JSON(w, res, http.StatusOK)
	}
}

func LoginHandler() api.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		if gothUser, err := gothic.CompleteUserAuth(w, r); err == nil {
			slog.Info("auth", "message", "user logged in successfully without redirect", "user", gothUser)
		} else {
			gothic.BeginAuthHandler(w, r)
		}
		return nil
	}
}
