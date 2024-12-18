package auth

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/kacperhemperek/twitter-v2/api"
	"github.com/kacperhemperek/twitter-v2/services"
)

type Middleware = func(h api.HandlerFunc) api.HandlerFunc

func NewAuthMiddleware(userService services.UserService, sessionService SessionService) Middleware {
	return func(h api.HandlerFunc) api.HandlerFunc {
		return func(w http.ResponseWriter, r *api.Request) (err error) {
			defer func() {
				var apiErr *api.APIError
				if err != nil && errors.As(err, &apiErr) && apiErr.Status == http.StatusUnauthorized {
					slog.Info("auth middleware", "message", "user is unauthorized, clearing session cookie")
					ClearSessionCookie(w)
				}
			}()
			sessCookie, err := r.Cookie("sessionID")
			if err != nil {
				slog.Debug("auth middleware", "message", "session cookie not found")
				return api.NewUnauthorizedError()
			}
			sess, err := sessionService.GetSession(r.Context(), sessCookie.Value)

			if err != nil && errors.Is(err, ErrSessionNotFound) {
				slog.Debug("auth middleware", "message", "session not found")
				return api.NewUnauthorizedError()
			}

			if err != nil {
				slog.Debug("auth middleware", "message", "unexpected error while getting session from store", "error", err)
				return err
			}

			if sess.IsExpired() {
				slog.Debug("auth middleware", "message", "session is expired")
				return api.NewUnauthorizedError()
			}

			user, err := userService.GetByID(r.Context(), sess.UserID)
			if err != nil && errors.Is(err, services.ErrUserNotFound) {
				slog.Debug("auth middleware", "message", "user not found")
				return api.NewUnauthorizedError()
			}
			if err != nil {
				slog.Debug("auth middleware", "message", "unexpected error while getting user from store", "error", err)
				return err
			}

			r.SetUser(user)
			r.SetSession(sess)
			slog.Debug("auth middleware", "message", "user is authorized, session is setup", "user", user)

			return h(w, r)
		}
	}
}
