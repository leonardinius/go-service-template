package apiworker

import (
	"bufio"
	"io"
	"net"
	"net/http"
)

type stdResponseWriter struct {
	io.Writer
	header     http.Header
	statusCode int
}

var (
	_ http.ResponseWriter = (*stdResponseWriter)(nil)
	_ http.Flusher        = (*stdResponseWriter)(nil)
	_ http.Hijacker       = (*stdResponseWriter)(nil)
)

// Header implements http.ResponseWriter.
func (w *stdResponseWriter) Header() http.Header {
	return w.header
}

// WriteHeader implements http.ResponseWriter.
func (w *stdResponseWriter) WriteHeader(status int) {
	w.statusCode = status
}

// Flush implements http.Flusher.
func (w *stdResponseWriter) Flush() {
	z := w.Writer
	if f, ok := z.(http.Flusher); ok {
		f.Flush()
	}
}

// Hijack implements http.Hijacker.
func (w *stdResponseWriter) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	z := w.Writer
	if h, ok := z.(http.Hijacker); ok {
		return h.Hijack()
	}
	return nil, nil, http.ErrNotSupported
}

// Write implements http.ResponseWriter.
func (w *stdResponseWriter) Write(b []byte) (int, error) {
	return w.Writer.Write(b)
}

func NewStdResponseWriter(w io.Writer) *stdResponseWriter {
	return &stdResponseWriter{
		Writer: w,
		header: make(http.Header),
	}
}
