package api

import (
	"fmt"
	"log/slog"
)

func SetupLogger() {
	var lvl slog.Level
	if ENV.IsProd() {
		lvl = slog.LevelWarn
	} else if ENV.IsDebug() {
		lvl = slog.LevelDebug
	} else {
		lvl = slog.LevelInfo
	}
	slog.SetLogLoggerLevel(lvl)
}

func LogServiceError(service, method string, err error) {
	m := fmt.Sprintf("failed %s method", method)
	s := fmt.Sprintf("%s service", service)
	slog.Error(s, "message", m, "error", err)
}
