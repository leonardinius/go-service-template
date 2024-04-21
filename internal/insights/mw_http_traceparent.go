package insights

import (
	"net/http"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
)

type TraceparentHandlerMiddleware struct {
	next  http.Handler
	props propagation.TextMapPropagator
}

var _ http.Handler = (*TraceparentHandlerMiddleware)(nil)

// NewTraceparentHandlerMiddleware creates a new TraceparentHandlerMiddleware.
// It injects the "traceparent" header into the response headers according
// to w3c trace context specification.
func NewTraceparentHandlerMiddleware(next http.Handler) http.Handler {
	return &TraceparentHandlerMiddleware{
		next:  next,
		props: otel.GetTextMapPropagator(),
	}
}

func (h *TraceparentHandlerMiddleware) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	h.props.Inject(req.Context(), propagation.HeaderCarrier(w.Header()))
	h.next.ServeHTTP(w, req)
}
