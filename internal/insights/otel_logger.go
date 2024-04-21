package insights

import (
	"log/slog"

	slogotel "github.com/remychantenay/slog-otel"
)

// NewLogOtelMiddleware wraps the provided slog.Handler with OpenTelemetry logging middleware.
func NewLogOtelMiddleware(next slog.Handler) slog.Handler {
	return slogotel.New(next)
}
