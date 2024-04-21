package apiserv

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"runtime"
)

func NewRecoveryHandlerMiddleware(handler http.Handler, logger *slog.Logger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func(ctx context.Context) {
			if err := recover(); err != nil {
				err := newPanickedErrorFromStack(err)
				logger.LogAttrs(ctx, slog.LevelError, "recovered from panic",
					slog.String("error", err.Error()),
					slog.String("stack", err.Stack()),
				)
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			}
		}(r.Context())
		handler.ServeHTTP(w, r)
	})
}

func newPanickedErrorFromStack(p any) *PanicError {
	stack := make([]byte, 64<<10)
	stack = stack[:runtime.Stack(stack, false)]
	return &PanicError{panic: p, stack: stack}
}

type PanicError struct {
	panic any
	stack []byte
}

func (e *PanicError) Error() string {
	return fmt.Sprintf("panic caught: %v", e.panic)
}

func (e *PanicError) Stack() string {
	return string(e.stack)
}
