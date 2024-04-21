package apiserv

import (
	"bufio"
	"context"
	"log/slog"
	"net"
	"net/http"
	"strings"
	"time"
)

func NewLogHandlerMiddleware(handler http.Handler, logger *slog.Logger, level slog.Level, message string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		started := time.Now()
		interceptWriter := stdResponseWriter{w, 0, 0}

		defer func(ctx context.Context) {
			took := time.Since(started)
			remoteAddr := requestGetRemoteAddress(r)
			requestLogLevel := level
			if interceptWriter.aHTTPStatus >= 400 {
				requestLogLevel = slog.LevelError
			}

			logger.WithGroup("http").
				LogAttrs(ctx, requestLogLevel, message,
					slog.String("remote_addr", remoteAddr),
					slog.String("proto", r.Proto),
					slog.String("method", r.Method),
					slog.String("path", r.URL.Path),
					slog.String("query", r.URL.RawQuery),
					slog.String("user_agent", r.UserAgent()),
					slog.Int64("content_length", r.ContentLength),
					slog.Int("status_code", interceptWriter.aHTTPStatus),
					slog.Int("response_bytes", interceptWriter.aResponseSizeBytes),
					slog.Duration("duration", took))
		}(r.Context())

		handler.ServeHTTP(&interceptWriter, r)
	})
}

type stdResponseWriter struct {
	http.ResponseWriter
	aHTTPStatus        int
	aResponseSizeBytes int
}

var (
	_ http.ResponseWriter = (*stdResponseWriter)(nil)
	_ http.Flusher        = (*stdResponseWriter)(nil)
	_ http.Hijacker       = (*stdResponseWriter)(nil)
)

// WriteHeader implements http.ResponseWriter.
func (w *stdResponseWriter) WriteHeader(status int) {
	w.aHTTPStatus = status
	w.ResponseWriter.WriteHeader(status)
}

// Flush implements http.Flusher.
func (w *stdResponseWriter) Flush() {
	z := w.ResponseWriter
	if f, ok := z.(http.Flusher); ok {
		f.Flush()
	}
}

// Hijack implements http.Hijacker.
func (w *stdResponseWriter) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	z := w.ResponseWriter
	if h, ok := z.(http.Hijacker); ok {
		return h.Hijack()
	}
	return nil, nil, http.ErrNotSupported
}

// Write implements http.ResponseWriter.
func (w *stdResponseWriter) Write(b []byte) (bytesWritten int, err error) {
	if w.aHTTPStatus == 0 {
		w.aHTTPStatus = 200
	}
	bytesWritten, err = w.ResponseWriter.Write(b)
	w.aResponseSizeBytes += bytesWritten
	return
}

// Request.RemoteAddress contains port, which we want to remove
// i.e.: "[::1]:58292" => "[::1]".
func ipAddrFromRemoteAddr(s string) string {
	idx := strings.LastIndex(s, ":")
	if idx == -1 {
		return s
	}
	return s[:idx]
}

// requestGetRemoteAddress returns ip address of the client making the request,
// taking into account http proxies.
func requestGetRemoteAddress(r *http.Request) string {
	hdrRealIP := r.Header.Get("X-Real-Ip")
	hdrForwardedFor := r.Header.Get("X-Forwarded-For")
	if hdrRealIP == "" && hdrForwardedFor == "" {
		return ipAddrFromRemoteAddr(r.RemoteAddr)
	}
	if hdrForwardedFor != "" {
		// X-Forwarded-For is potentially a list of addresses separated with ","
		parts := strings.Split(hdrForwardedFor, ",")
		for i, p := range parts {
			parts[i] = strings.TrimSpace(p)
		}
		return parts[0]
	}
	return hdrRealIP
}
