package cmd

import (
	"context"
	"log/slog"
	"os"
	"strings"
	"sync"

	"github.com/leonardinius/go-service-template/internal/log"
)

const logDefaultLevel = slog.LevelInfo

var logLevelNamesMapping = map[string]slog.Level{
	"debug": slog.LevelDebug,
	"info":  slog.LevelInfo,
	"warn":  slog.LevelWarn,
	"error": slog.LevelError,
}

var logInitOnce sync.Once = sync.Once{}

func MustSetupLogger(ctx context.Context, level string) {
	// Initialize logger
	// Do this only once, because of parallel e2e tests.
	logInitOnce.Do(func() {
		logLevel := logDefaultLevel
		if levelLookup, ok := logLevelNamesMapping[strings.ToLower(level)]; ok {
			logLevel = levelLookup
		}

		logger := log.InitDefaultLogger(os.Stdout, logLevel)
		log.SetOtelGlobalLogger(ctx, logger)
	})
}
