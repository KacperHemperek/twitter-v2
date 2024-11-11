package auth

import (
	"errors"
	"fmt"
	"log/slog"
	"math/rand/v2"
	"net/http"
	"strings"

	"github.com/kacperhemperek/twitter-v2/api"
	"github.com/kacperhemperek/twitter-v2/services"
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

func AuthCallbackHanlder(userService services.UserService) api.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		// TODO: i think should set something that tells frontend that user is
		// authenticated and that api can then read on request to get user info

		gothUser, err := gothic.CompleteUserAuth(w, r)

		if err != nil {
			return err
		}

		_, err = userService.GetByEmail(r.Context(), gothUser.Email)
		isNotFoundError := errors.Is(err, services.ErrUserNotFound)

		if err != nil && !isNotFoundError {
			return err
		}

		if err != nil && isNotFoundError {
			// NOTE: some providers (google included) do not return name and
			// instead they split it to first name and last name
			// first name and last name are missing as well just take first part of email
			name := strings.TrimSpace(gothUser.Name)
			if len(name) == 0 {
				name = strings.TrimSpace(fmt.Sprintf("%s %s", gothUser.FirstName, gothUser.LastName))
			}
			if len(name) == 0 {
				name = strings.Split(gothUser.Email, "@")[0]
			}
			if len(name) == 0 {
				name = generateRandomUserName()
			}
			_, err = userService.CreateUser(r.Context(), gothUser.Email, name)
			if err != nil {
				return err
			}
		}

		http.Redirect(w, r, api.ENV.FRONTEND_URL, http.StatusTemporaryRedirect)
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
			// TODO: if user completed logging in get user from the database
			slog.Info("auth", "message", "user logged in successfully without redirect", "user", gothUser)
		} else {
			gothic.BeginAuthHandler(w, r)
		}
		return nil
	}
}

func generateRandomUserName() string {
	return fmt.Sprintf("User%d", rand.IntN(99999-10000+1)+10000)
}
