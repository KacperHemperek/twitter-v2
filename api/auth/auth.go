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

// TODO:
// - replace dev google keys with prod keys on prod server
// - add more options than just google to auth (discord, twitter(?))

func AuthCallbackHanlder(userService services.UserService) api.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		gothUser, err := gothic.CompleteUserAuth(w, r)

		if err != nil {
			return err
		}

		user, err := getOrCreateUser(r.Context(), gothUser, userService)

		// TODO: store access and refresh token in query params so the login page on frontend
		// can store it in the memory
		redirectUrlStr, err := url.JoinPath(api.ENV.FRONTEND_URL, "login", "success")
		if err != nil {
			return err
		}
		redirectUrl, err := url.Parse(redirectUrlStr)
		if err != nil {
			return err
		}
		token := NewUserToken(user.ID, user.Email, user.Name)
		accessToken, err := token.SignAccessToken()
		if err != nil {
			return err
		}
		refreshToken, err := token.SignRefreshToken()
		if err != nil {
			return err
		}
		query := redirectUrl.Query()
		query.Add("accessToken", accessToken)
		query.Add("refreshToken", refreshToken)
		redirectUrl.RawQuery = query.Encode()

		http.Redirect(w, r, redirectUrl.String(), http.StatusTemporaryRedirect)
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

func LoginHandler(userService services.UserService) api.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		if gothUser, err := gothic.CompleteUserAuth(w, r); err == nil {
			user, nil := getOrCreateUser(r.Context(), gothUser, userService)

			if err != nil {
				return nil
			}

			redirectUrlStr, err := url.JoinPath(api.ENV.FRONTEND_URL, "login", "success")
			if err != nil {
				return err
			}
			redirectUrl, err := url.Parse(redirectUrlStr)
			if err != nil {
				return err
			}
			token := NewUserToken(user.ID, user.Email, user.Name)
			accessToken, err := token.SignAccessToken()
			if err != nil {
				return err
			}
			refreshToken, err := token.SignRefreshToken()
			if err != nil {
				return err
			}
			query := redirectUrl.Query()
			query.Add("accessToken", accessToken)
			query.Add("refreshToken", refreshToken)
			redirectUrl.RawQuery = query.Encode()

			http.Redirect(w, r, redirectUrl.String(), http.StatusTemporaryRedirect)
			// TODO: parse token and add to client response
			slog.Info("auth", "message", "user logged in successfully without redirect", "user", gothUser)
		} else {
			gothic.BeginAuthHandler(w, r)
		}
		return nil
	}
}

func GetMeHandler() api.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		return api.JSON(w, map[string]string{"user": "user info"}, http.StatusOK)
	}
}

func getOrCreateUser(ctx context.Context, gothUser goth.User, userService services.UserService) (user *services.UserModel, err error) {

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
