package auth

import (
	"log/slog"
	"net/http"

	"github.com/kacperhemperek/twitter-v2/api"
)

type Middleware = func(h api.HandlerFunc) api.HandlerFunc

func NewAuthMiddleware() Middleware {
	return func(h api.HandlerFunc) api.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) error {
			// TODO: add checks for tokens and update access token when refresh token is valid
			accessTokenStr := r.Header.Get("X-Access-Token")
			// refreshTokenStr := r.Header.Get("X-Refresh-Token")
			accessToken, err := ParseToken(accessTokenStr)

			if err != nil {
				return err
			}

			slog.Info("auth middleware", "userToken", accessToken)

			return h(w, r)
		}
	}
}
