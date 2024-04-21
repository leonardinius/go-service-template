package log

import (
	"context"
	"log/slog"

	"github.com/go-logr/logr"
	"go.opentelemetry.io/otel"
)

func SetOtelGlobalLogger(ctx context.Context, logger *slog.Logger) {
	logrLogger := logr.FromSlogHandler(logger.Handler())
	otel.SetLogger(logrLogger)
	otel.SetErrorHandler(otel.ErrorHandlerFunc(func(err error) {
		logger.LogAttrs(ctx,
			slog.LevelError,
			"OpenTelemetry error",
			slog.String("error", err.Error()))
	}))
}
