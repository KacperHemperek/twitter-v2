package api

import (
	"fmt"
	"log/slog"
)

func LogServiceError(service, method string, err error) {
	m := fmt.Sprintf("failed %s method", method)
	s := fmt.Sprintf("%s service", service)
	slog.Error(s, "message", m, "error", err)
}
