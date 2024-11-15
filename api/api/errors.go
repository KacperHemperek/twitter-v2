package api

import "net/http"

func NewBadRequestError(m string) *APIError {
	return &APIError{
		Message: m,
		Status:  http.StatusUnauthorized,
	}
}

func NewUnauthorizedError() *APIError {
	return &APIError{
		Message: "user must be authenticated to access this endpoint",
		Status:  http.StatusUnauthorized,
	}
}
