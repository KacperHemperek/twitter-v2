package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
)

type HandlerFunc = func(w http.ResponseWriter, r *http.Request) error

func Handle(h HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := h(w, r)

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
