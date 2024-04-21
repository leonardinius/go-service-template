package apiserv

import (
	"context"
	"net"
	"net/http"
)

// ListenAndServe starts an HTTP server and listens for incoming requests.
// It takes a context.Context, an *servehttp.Server, and optional configuration options.
// The returned error indicates if there was an error starting the server or shutting it down.
// It uses OpenTelemetry to set up tracing and metrics for the server.
//
// The server’s BaseContext is set to the provided context.Context.
//
// If the context.Context is canceled, the server is gracefully shutdown.
//
// Example usage:
// err := ListenAndServe(ctx, srv)
//
//	if err != nil {
//	    log.Fatal(err)
//	}
func ListenAndServe(parentContext context.Context, srv *http.Server, options ...ListenOption) (err error) {
	ctx, stop := context.WithCancel(parentContext)
	defer stop()
	serverConfig := initializeListenOptions(options)

	// Start HTTP server and listen for incoming requests.
	srv.BaseContext = func(_ net.Listener) context.Context { return ctx }
	return serverConfig.ListenAndServe(srv)
}

type serveListenConfigOptions struct {
	isSsl             bool
	certFile, keyFile string
}

func (cfg *serveListenConfigOptions) ListenAndServe(srv *http.Server) (err error) {
	if cfg.isSsl {
		return srv.ListenAndServeTLS(cfg.certFile, cfg.keyFile)
	}

	return srv.ListenAndServe()
}

// ListenOption is an interface that represents a configuration option for the serveListenConfigOptions struct.
// All implementations of ListenOption must implement the apply method,
// which takes a *serveConfigOptions parameter and applies the configuration option to it.
type ListenOption interface {
	apply(option *serveListenConfigOptions)
}

type listenOptionFunc func(*serveListenConfigOptions)

func (f listenOptionFunc) apply(srv *serveListenConfigOptions) {
	f(srv)
}

func initializeListenOptions(options []ListenOption) *serveListenConfigOptions {
	// init serveConfig with default context and stop function
	serveConfig := &serveListenConfigOptions{}
	for _, option := range options {
		option.apply(serveConfig)
	}
	return serveConfig
}

// WithListenSSL returns an Option that configures the server with SSL.
//
// The certFile parameter specifies the path to the certificate file.
//
// The keyFile parameter specifies the path to the private key file.
//
// Example usage:
//
//	opts := []Option{
//	  WithListenSSL("cert.pem", "key.pem"),
//	}
//	server := NewServer(opts...)
//
// The server will use the provided SSL certificate and private key files
// to servehttp the connections securely over HTTPS.
// If WithListenSSL is not used, the server will listen on HTTP.
func WithListenSSL(certFile, keyFile string) ListenOption {
	return listenOptionFunc(func(srv *serveListenConfigOptions) {
		srv.isSsl = true
		srv.certFile = certFile
		srv.keyFile = keyFile
	})
}
