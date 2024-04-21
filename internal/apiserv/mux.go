package apiserv

import (
	"context"
	"log/slog"
	"net/http"
	"net/url"
	"strings"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"

	"github.com/leonardinius/go-service-template/internal/insights"
)

type Route struct {
	// GET /helloworld
	Pattern string
	Handler http.Handler
}

func NewRoute(pattern string, handler http.Handler) Route {
	return Route{Pattern: pattern, Handler: handler}
}

func (r *Route) Path() string {
	if parts := strings.SplitN(r.Pattern, " ", 2); len(parts) > 1 {
		return parts[1]
	}

	return r.Pattern
}

// BuildHTTPMux creates an HTTP handler based on the provided register function.
//
// The registerFn parameter is a function that registers handlers for specific patterns.
// It takes a RegisterHandlerFn as a parameter, which is a function that associates a pattern with a handler function.
//
// The BuildHTTPMux function creates a new ServeMux and uses the handleFunc helper function to register handlers using the registerFn.
//
// Example usage:
//
//	ctx, cancel := context.WithCancel(context.Background())
//	defer cancel()
//
//	handler := httpserv.BuildHTTPMux(ctx, func(register httpserv.RegisterHandlerFn) {
//		register("GET /helloworld", func(w servehttp.ResponseWriter, r *servehttp.Request) {
//			w.WriteHeader(http.StatusOK)
//			_, _ = w.Write([]byte("Hello World"))
//		})
//	})
//
//	srv := &servehttp.Server{
//		Addr:    address,
//		Handler: handler,
//	}
//
//	err := httpserv.ListenAndServe(srv, httpserv.WithContext(ctx, cancel))
//
//	if err != nil {
//		log.Fatal(err)
//	}`
func BuildHTTPMux(ctx context.Context, options ...Option) http.Handler {
	c := initializeOptions(options)
	address, routes, logger, level := c.address, c.routes, c.logger, c.middlewareLogLevel

	mux := http.NewServeMux()
	mux = registerHandlers(ctx, mux, address, routes)

	// Add HTTP instrumentation for the whole server.
	var handler http.Handler = mux
	// Please note the order of middleware registration is important.
	// Execution is the reverse of the registration order.
	handler = NewRecoveryHandlerMiddleware(handler, logger)
	handler = insights.NewAddXHeadersHandlerMiddleware(handler)
	handler = NewLogHandlerMiddleware(handler, logger, level, "http")
	handler = insights.NewTraceparentHandlerMiddleware(handler)
	handler = insights.NewOtelHandlerMiddleware(handler, "http")
	handler = insights.NewMetricsHandlerMiddleware(handler, ctx, "http")
	return handler
}

func registerHandlers(ctx context.Context, mux *http.ServeMux, address string, routes []Route) *http.ServeMux {
	// registerFn is a middleware that registers handlers for specific patterns,
	registerFn := func(pattern string, handler http.Handler) {
		handlerID := handlerIDFromPattern(pattern)
		handler = otelhttp.WithRouteTag(handlerID, handler)
		mux.Handle(pattern, handler)
		logHandlerRegistered(ctx, pattern, address)
	}

	for _, route := range routes {
		registerFn(route.Pattern, route.Handler)
	}

	return mux
}

func handlerIDFromPattern(pattern string) string {
	if parts := strings.SplitN(pattern, " ", 2); len(parts) > 1 {
		return parts[1]
	}

	return pattern
}

func logHandlerRegistered(ctx context.Context, pattern, address string) {
	method := "ALL"
	path := pattern

	if parts := strings.SplitN(path, " ", 2); len(parts) > 1 {
		method = strings.ToUpper(parts[0])
		path = parts[1]
	}

	registeredURL := &url.URL{
		Scheme: "http",
		Host:   address,
		Path:   path,
	}

	slog.LogAttrs(ctx, slog.LevelInfo, "http handler registered",
		slog.String("method", method),
		slog.String("url", registeredURL.String()))
}

type serverConfigOptions struct {
	logger             *slog.Logger
	middlewareLogLevel slog.Level
	address            string
	routes             []Route
}

type Option interface {
	apply(option *serverConfigOptions)
}

type optionFunc func(*serverConfigOptions)

func (f optionFunc) apply(srv *serverConfigOptions) {
	f(srv)
}

func initializeOptions(opts []Option) *serverConfigOptions {
	cfg := &serverConfigOptions{
		logger:             slog.Default(),
		middlewareLogLevel: slog.LevelInfo,
		address:            "",
	}

	for _, opt := range opts {
		opt.apply(cfg)
	}
	return cfg
}

// WithLogger returns an Option that configures the server with the given logger.
//
// The logger parameter specifies the logger to use for logging.
//
// Example usage:
//
//	opts := []Option{
//	  WithLogger(logger),
//	}
//	server := NewServer(opts...)
//
// The server will use the provided logger for logging.
func WithLogger(logger *slog.Logger) Option {
	return optionFunc(func(srv *serverConfigOptions) {
		srv.logger = logger
	})
}

// WithMiddlewareLogLevel returns an Option that configures the server with the given middleware log level.
//
// The level parameter specifies the log level to use for middleware logging.
//
// Example usage:
//
//	opts := []Option{
//	  WithMiddlewareLogLevel(slog.LevelInfo),
//	}
//	server := NewServer(opts...)
//
// The server will use the provided log level for middleware logging.
func WithMiddlewareLogLevel(level slog.Level) Option {
	return optionFunc(func(srv *serverConfigOptions) {
		srv.middlewareLogLevel = level
	})
}

// WithAddress returns an Option that configures the server with the given address.
//
// The address parameter specifies the address to use for the server.
//
// Example usage:
//
//	opts := []Option{
//	  WithAddress(address),
//	}
//	server := NewServer(opts...)
//
// The server will use the provided address for the server.
func WithAddress(address string) Option {
	return optionFunc(func(srv *serverConfigOptions) {
		srv.address = address
	})
}

// WithRoute returns an Option that configures the server with the given route.
//
// The pattern parameter specifies the pattern to use for the route.
//
// The handler parameter specifies the handler to use for the route.
//
// Example usage:
//
//	opts := []Option{
//	  WithRoute("/helloworld", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
//	    w.WriteHeader(http.StatusOK)
//	    _, _ = w.Write([]byte("Hello World"))
//	  })),
//	}
//	server := NewServer(opts...)
//
// The server will use the provided route for the server.
func WithRoute(pattern string, handler http.Handler) Option {
	return optionFunc(func(srv *serverConfigOptions) {
		srv.routes = append(srv.routes, Route{Pattern: pattern, Handler: handler})
	})
}

// WithRoutes returns an Option that configures the server with the given routes.
//
// The routes parameter specifies the routes to use for the server.
//
// Example usage:
//
//	opts := []Option{
//	  WithRoutes(routes...),
//	}
//	server := NewServer(opts...)
//
// The server will use the provided routes for the server.
func WithRoutes(routes ...Route) Option {
	return optionFunc(func(srv *serverConfigOptions) {
		srv.routes = append(srv.routes, routes...)
	})
}
