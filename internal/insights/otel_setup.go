package insights

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"

	"github.com/leonardinius/go-service-template/internal/services/version"

	semconv "go.opentelemetry.io/otel/semconv/v1.24.0"
)

var (
	// Default is 5 seconds.
	traceBatchTimeout = 5 * time.Second

	errUnsupportedOLTPProtocol = errors.New("unsupported otlp protocol, supported protocols are grpc, http/protobuf")
)

// SetupOtelSDK bootstraps the OpenTelemetry pipeline.
// If it does not return an error, make sure to call shutdown for proper cleanup.
func SetupOtelSDK(ctx context.Context) (func(context.Context) error, error) {
	var shutdownFuncs []func(context.Context) error

	// shutdown calls cleanup functions registered via shutdownFuncs.
	// The errors from the calls are joined.
	// Each registered cleanup will be invoked once.
	shutdown := func(ctx context.Context) error {
		var shutdownErr error
		for _, fn := range shutdownFuncs {
			shutdownErr = errors.Join(shutdownErr, fn(ctx))
		}
		shutdownFuncs = nil
		return shutdownErr
	}

	// handleErr calls shutdown for cleanup and makes sure that all errors are returned.
	handleErr := func(inErr error) (func(context.Context) error, error) {
		return shutdown, errors.Join(inErr, shutdown(ctx))
	}

	// Set up propagator.
	prop := newTextMapPropagator()
	otel.SetTextMapPropagator(prop)

	// Set up trace provider.
	tracerProvider, err := newTraceProvider(ctx)
	if err != nil {
		return handleErr(err)
	}
	otel.SetTracerProvider(tracerProvider)
	shutdownFuncs = append(shutdownFuncs, tracerProvider.Shutdown)

	// Set up meter provider.
	meterProvider, err := newMeterProvider()
	if err != nil {
		return handleErr(err)
	}
	shutdownFuncs = append(shutdownFuncs, meterProvider.Shutdown)
	otel.SetMeterProvider(meterProvider)

	return shutdown, nil
}

func newTextMapPropagator() propagation.TextMapPropagator {
	return propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	)
}

func newTraceProvider(ctx context.Context) (*trace.TracerProvider, error) {
	traceExporter, err := newTraceSpanExporter(ctx)
	if err != nil {
		return nil, err
	}

	res, err := newServiceResource(ctx)
	if err != nil {
		return nil, err
	}

	traceProvider := trace.NewTracerProvider(
		trace.WithBatcher(traceExporter, trace.WithBatchTimeout(traceBatchTimeout)),
		trace.WithResource(res),
	)
	return traceProvider, nil
}

func newTraceSpanExporter(ctx context.Context) (trace.SpanExporter, error) {
	set := os.Getenv("OTEL_TRACES_EXPORTER") == "otlp" ||
		os.Getenv("OTEL_EXPORTER_OTLP_ENDPOINT") != "" ||
		os.Getenv("OTEL_EXPORTER_OTLP_TRACES_ENDPOINT") != ""
	if !set {
		return nil, nil
	}

	proto := os.Getenv("OTEL_EXPORTER_OTLP_TRACES_PROTOCOL")
	if proto == "" {
		proto = os.Getenv("OTEL_EXPORTER_OTLP_PROTOCOL")
	}
	if proto == "" {
		proto = "grpc"
	}

	var c otlptrace.Client

	switch proto {
	case "grpc":
		c = otlptracegrpc.NewClient()
	case "http/protobuf":
		c = otlptracehttp.NewClient()
	// case "http/json": // unsupported by library
	default:
		return nil, wrapSupportedOLTPProtocol(proto)
	}

	return otlptrace.New(ctx, c)
}

//nolint:unparam // error is always nil.
func newMeterProvider() (*metric.MeterProvider, error) {
	return metric.NewMeterProvider(), nil
}

func newServiceResource(ctx context.Context) (*resource.Resource, error) {
	res, err := resource.New(ctx,
		// Add custom resource attributes.
		// Note: will be overridden by the same attributes from environment variables if present.
		resource.WithAttributes(
			semconv.ServiceName(version.ServiceName),
			semconv.ServiceVersion(version.FullVersion),
		),
		resource.WithFromEnv(),      // Discover and provide attributes from OTEL_RESOURCE_ATTRIBUTES and OTEL_SERVICE_NAME env variables.
		resource.WithTelemetrySDK(), // Discover and provide information about the OpenTelemetry SDK used.
		resource.WithProcess(),      // Discover and provide process information.
		resource.WithOS(),           // Discover and provide OS information.
		resource.WithContainer(),    // Discover and provide container information.
		resource.WithHost(),         // Discover and provide host information.
	)

	if errors.Is(err, resource.ErrPartialResource) || errors.Is(err, resource.ErrSchemaURLConflict) {
		slog.LogAttrs(ctx, slog.LevelWarn, "resource creation failed", slog.String("error", err.Error()))
	} else if err != nil {
		return nil, err
	}

	return res, nil
}

func wrapSupportedOLTPProtocol(protocol string) error {
	return fmt.Errorf("%w: %s", errUnsupportedOLTPProtocol, protocol)
}
