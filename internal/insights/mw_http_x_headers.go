package insights

import (
	"net/http"

	"go.opentelemetry.io/otel/trace"
)

type AddXHeadersMiddleware struct {
	next http.Handler
}

var _ http.Handler = (*AddXHeadersMiddleware)(nil)

// NewAddXHeadersHandlerMiddleware creates a new AddXHeadersMiddleware.
// It adds the trace ID as rthe "X-Trace-Id" response headers if a trace is active.
func NewAddXHeadersHandlerMiddleware(next http.Handler) http.Handler {
	return &AddXHeadersMiddleware{
		next: next,
	}
}

func (h *AddXHeadersMiddleware) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	if s := trace.SpanContextFromContext(req.Context()); s.IsValid() {
		w.Header().Set("X-Trace-Id", s.TraceID().String())
	}
	h.next.ServeHTTP(w, req)
}
