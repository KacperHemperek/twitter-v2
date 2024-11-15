package auth

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"math/rand/v2"
	"net/http"
	"net/url"
	"strings"

	"github.com/kacperhemperek/twitter-v2/api"
	"github.com/kacperhemperek/twitter-v2/models"
	"github.com/kacperhemperek/twitter-v2/services"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/google"
)

func Setup() {
	slog.Info("auth", "message", "setting up oauth")
	callbackURL, err := url.JoinPath(api.ENV.API_URL, "api", "auth", "google", "callback")
	if err != nil {
		slog.Error("auth", "message", "could not create callback url", "error", err)
	}
	goth.UseProviders(
		google.New(api.ENV.GOOGLE_CLIENT_ID, api.ENV.GOOGLE_CLIENT_SECRET, callbackURL),
	)
	slog.Info("auth", "message", "oauth correctly setup")
}

func AuthCallbackHanlder(userService services.UserService, sessionService SessionService) api.HandlerFunc {
	return func(w http.ResponseWriter, r *api.Request) error {
		user, err := r.User()

		if err == nil && user != nil {
			return api.NewBadRequestError(
				"user is already logged in")
		}

		gothUser, err := gothic.CompleteUserAuth(w, r.Request)

		if err != nil {
			return err
		}

		user, err = getOrCreateUser(r.Context(), gothUser, userService)

		redirectUrlStr, err := url.JoinPath(api.ENV.FRONTEND_URL, "login", "success")
		if err != nil {
			return err
		}
		sess, err := sessionService.CreateSession(r.Context(), user.ID)
		if err != nil {
			return err
		}

		SetSessionCookie(w, sess.ID)
		http.Redirect(w, r.Request, redirectUrlStr, http.StatusTemporaryRedirect)
		return nil
	}
}

func LogoutHandler(sessionService SessionService) api.HandlerFunc {
	type response struct {
		Message string `json:"message"`
	}
	return func(w http.ResponseWriter, r *api.Request) error {
		sess, err := r.Session()
		if err != nil {
			return err
		}
		err = sessionService.DeleteSession(r.Context(), sess.ID)
		if err != nil {
			return err
		}
		gothic.Logout(w, r.Request)
		ClearSessionCookie(w)
		res := &response{
			Message: "user logged out successfully",
		}
		return api.JSON(w, res, http.StatusOK)
	}
}

func LoginHandler(userService services.UserService, sessionService SessionService) api.HandlerFunc {
	return func(w http.ResponseWriter, r *api.Request) error {
		user, err := r.User()

		if err == nil && user != nil {
			return &api.APIError{
				Message: "user is already logged in",
				Status:  http.StatusBadRequest,
			}
		}

		if gothUser, err := gothic.CompleteUserAuth(w, r.Request); err == nil {
			user, nil := getOrCreateUser(r.Context(), gothUser, userService)

			if err != nil {
				return nil
			}

			redirectUrlStr, err := url.JoinPath(api.ENV.FRONTEND_URL, "login", "success")
			if err != nil {
				return err
			}
			sess, err := sessionService.CreateSession(r.Context(), user.ID)
			if err != nil {
				return err
			}
			SetSessionCookie(w, sess.ID)
			http.Redirect(w, r.Request, redirectUrlStr, http.StatusTemporaryRedirect)
		} else {
			gothic.BeginAuthHandler(w, r.Request)
		}
		return nil
	}
}

func GetMeHandler() api.HandlerFunc {
	return func(w http.ResponseWriter, r *api.Request) error {
		user, err := r.User()
		if err != nil {
			return err
		}
		return api.JSON(w, map[string]any{"user": user}, http.StatusOK)
	}
}

func getOrCreateUser(ctx context.Context, gothUser goth.User, userService services.UserService) (user *models.UserModel, err error) {

	user, err = userService.GetByEmail(ctx, gothUser.Email)
	isNotFoundError := errors.Is(err, services.ErrUserNotFound)

	if err != nil && !isNotFoundError {
		return nil, err
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
		user, err = userService.CreateUser(
			ctx,
			gothUser.Email,
			name,
			gothUser.AvatarURL,
		)
		if err != nil {
			return nil, err
		}
	}

	return user, nil
}

func generateRandomUserName() string {
	return fmt.Sprintf("User%d", rand.IntN(99999-10000+1)+10000)
}
