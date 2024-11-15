package api

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/kacperhemperek/twitter-v2/models"
)

var (
	ErrNoUserInRequest    = errors.New("user not found in request context")
	ErrNoSessionInRequest = errors.New("session not found in request context")
)

const SESSION_CTX_KEY = "session"
const USER_CTX_KEY = "user"

type Request struct {
	*http.Request
}

func (r *Request) User() (*models.UserModel, error) {
	user, ok := r.Context().Value(USER_CTX_KEY).(*models.UserModel)
	if !ok {
		return nil, ErrNoUserInRequest
	}

	return user, nil
}

func (r *Request) SetUser(u *models.UserModel) {
	ctx := context.WithValue(r.Context(), USER_CTX_KEY, u)
	r.Request = r.WithContext(ctx)
}

func (r *Request) Session() (*models.SessionModel, error) {
	session, ok := r.Context().Value(SESSION_CTX_KEY).(*models.SessionModel)
	if !ok {
		return nil, ErrNoSessionInRequest
	}

	return session, nil
}

func (r *Request) SetSession(s *models.SessionModel) {
	ctx := context.WithValue(r.Context(), SESSION_CTX_KEY, s)
	r.Request = r.WithContext(ctx)
}

type HandlerFunc = func(w http.ResponseWriter, r *Request) error

func Handle(h HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cr := &Request{
			Request: r,
		}
		err := h(w, cr)

		if err != nil {
			var apiError *APIError
			if errors.As(err, &apiError) {
				slog.Error(
					"api http handler",
					"method", r.Method,
					"url", r.URL.Path,
					"error", apiError.Error(),
					"status", apiError.Status,
				)

				jsonErr := JSON(w, apiError, apiError.Status)
				if jsonErr != nil {
					slog.Error("json return error", "error", jsonErr.Error())
					// This should never happen if it does program should panic
					panic(1)
				}
			} else {
				slog.Error(
					"unhandled http handler",
					"method", r.Method,
					"url", r.URL.Path,
					"error", err.Error(),
					"status", http.StatusInternalServerError,
				)
				apiError = NewAPIError("Unhandled Server Error", http.StatusInternalServerError)
				jsonErr := JSON(w, apiError, apiError.Status)
				if jsonErr != nil {
					slog.Error("json return error", "error", jsonErr.Error())
					// This should never happen if it does program should panic
					panic(1)
				}
			}
		} else {
			slog.Info(
				"request",
				"method", r.Method,
				"url", r.URL.Path,
			)
		}
	}
}

type APIError struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

func (e APIError) Error() string {
	return fmt.Sprintf("[%d] %s", e.Status, e.Message)
}

func NewAPIError(m string, s int) *APIError {
	return &APIError{
		Status:  s,
		Message: m,
	}
}

func JSON(w http.ResponseWriter, data any, status int) error {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(data)
}
