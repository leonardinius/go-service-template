package insights

import (
	"fmt"
	"net/http"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

// NewOtelHandlerMiddleware returns an HTTP handler with OpenTelemetry instrumentation.
// The handler parameter is the original HTTP handler to be instrumented.
// The second parameter is unused.
// It adds middleware to handle tracing and span name formatting.
// The span name formatter uses the HTTP method and URL path.
func NewOtelHandlerMiddleware(handler http.Handler, operation string) http.Handler {
	return otelhttp.NewHandler(handler, operation,
		otelhttp.WithPublicEndpoint(),
		otelhttp.WithSpanNameFormatter(func(operation string, r *http.Request) string {
			return fmt.Sprintf("%s %s %s", operation, r.Method, r.URL.Path)
		}))
}
