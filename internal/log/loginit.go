package log

import (
	"io"
	"log/slog"

	"github.com/leonardinius/go-service-template/internal/insights"
)

// InitDefaultLogger initializes a default logger with the default writer and log level.
func InitDefaultLogger(w io.Writer, level slog.Level) *slog.Logger {
	handler := NewJSONHandler(w, level)
	handler = insights.NewLogOtelMiddleware(handler)
	logger := NewLogger(handler)
	slog.SetLogLoggerLevel(level)
	slog.SetDefault(logger)
	return logger
}

func NewJSONHandler(w io.Writer, level slog.Level) slog.Handler {
	return slog.NewJSONHandler(w, &slog.HandlerOptions{
		Level: level,
	})
}

// NewLogger creates a new logger with the given handler.
func NewLogger(handler slog.Handler) *slog.Logger {
	return slog.New(handler)
}
