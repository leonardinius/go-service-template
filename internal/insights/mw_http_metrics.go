package insights

import (
	"context"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/slok/go-http-metrics/middleware"
	"github.com/slok/go-http-metrics/middleware/std"

	metrics "github.com/slok/go-http-metrics/metrics/prometheus"
)

func NewMetricsHandlerMiddleware(next http.Handler, ctx context.Context, service string) http.Handler {
	insightsMiddleware := middleware.New(middleware.Config{
		Service: service,
		Recorder: metrics.NewRecorder(metrics.Config{
			Registry: RegistrerFromContext(ctx),
		}),
	})

	return std.Handler("", insightsMiddleware, next)
}

type (
	registrerKey struct{}
	gathererKey  struct{}
)

func NewMetricsRegistry() *prometheus.Registry {
	reg := prometheus.NewRegistry()
	reg.MustRegister(collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}))
	reg.MustRegister(collectors.NewGoCollector())
	return reg
}

// ContextWithRegistry returns a new context with the given registry.
func ContextWithRegistry(ctx context.Context, registry *prometheus.Registry) context.Context {
	return ContextWithRegistrer(ContextWithGatherer(ctx, registry), registry)
}

// ContextWithRegistrer returns a new context with the given registry.
func ContextWithRegistrer(ctx context.Context, registrer prometheus.Registerer) context.Context {
	return context.WithValue(ctx, registrerKey{}, registrer)
}

// ContextWithGatherer returns a new context with the given registry.
func ContextWithGatherer(ctx context.Context, gatherer prometheus.Gatherer) context.Context {
	return context.WithValue(ctx, gathererKey{}, gatherer)
}

// RegistrerFromContext returns the registrer from the given context.
func RegistrerFromContext(ctx context.Context) prometheus.Registerer {
	v := ctx.Value(registrerKey{})

	if r, ok := v.(prometheus.Registerer); ok {
		return r
	}

	return prometheus.DefaultRegisterer
}

// GathererFromContext returns the gatherer from the given context.
func GathererFromContext(ctx context.Context) prometheus.Gatherer {
	v := ctx.Value(registrerKey{})

	if g, ok := v.(prometheus.Gatherer); ok {
		return g
	}

	return prometheus.DefaultGatherer
}

// NewMetricsHTTPHandler returns a new HTTP handler that exposes the metrics.
func NewMetricsHTTPHandler(ctx context.Context) http.Handler {
	registerer := RegistrerFromContext(ctx)
	gatherer := GathererFromContext(ctx)

	return promhttp.HandlerFor(gatherer, promhttp.HandlerOpts{Registry: registerer})
}
