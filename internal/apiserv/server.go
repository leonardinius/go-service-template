package apiserv

import (
	"context"
	"log/slog"
	"net/http"
	"time"

	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"

	"github.com/leonardinius/go-service-template/internal/insights"
)

const (
	HeaderReadTimeout = 3 * time.Second
	MetricsRoutePath  = "/metrics"
)

func NewServer(ctx context.Context, address string, options ...Option) (*http.Server, error) {
	options = append(options, WithAddress(address))
	handler := BuildHTTPMux(ctx, options...)

	srv := http.Server{
		Addr: address,
		// // Use h2c so we can serve HTTP/2 without TLS.
		Handler:           h2c.NewHandler(handler, &http2.Server{}),
		ReadHeaderTimeout: HeaderReadTimeout,
	}

	return &srv, nil
}

func NewDefaultServer(ctx context.Context, address string, routes ...Route) (*http.Server, error) {
	return NewServer(ctx, address,
		WithLogger(slog.Default()),
		WithMiddlewareLogLevel(slog.LevelDebug),
		WithRoutes(routes...),
		WithRoute("GET "+MetricsRoutePath, insights.NewMetricsHTTPHandler(ctx)),
	)
}
