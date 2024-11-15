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
				return api.NewUnauthorizedError()
			}
			sess, err := sessionService.GetSession(r.Context(), sessCookie.Value)

			if err != nil && errors.Is(err, ErrSessionNotFound) {
				return api.NewUnauthorizedError()
			}

			if err != nil {
				return err
			}

			if sess.IsExpired() {
				return api.NewUnauthorizedError()
			}

			user, err := userService.GetByID(r.Context(), sess.UserID)
			if err != nil && errors.Is(err, services.ErrUserNotFound) {
				return api.NewUnauthorizedError()
			}
			if err != nil {
				return err
			}

			r.SetUser(user)
			r.SetSession(sess)

			return h(w, r)
		}
	}
}
