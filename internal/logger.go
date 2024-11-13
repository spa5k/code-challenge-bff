package internal

import (
	"context"
	"log/slog"
	"os"

	"github.com/lmittmann/tint"
)

var loggerCache *slog.Logger

func NewLogger() *slog.Logger {
	w := os.Stderr
	logger := slog.New(tint.NewHandler(w, nil))
	loggerCache = logger

	return logger
}

func GetLogger(ctx context.Context) *slog.Logger {
	return loggerCache
}
