package apiserv_test

import (
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/leonardinius/go-service-template/internal/apiserv"
	"github.com/leonardinius/go-service-template/internal/log"
)

func TestNewLogHandlerMiddlewareSmokeTest(t *testing.T) {
	t.Parallel()
	// arrange: an handler that returns "Hello World!"
	helloWorld := "<html><body>Hello World!</body></html>"
	handler := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		_, _ = io.WriteString(w, helloWorld)
	})
	buffer := &strings.Builder{}
	// arrange: a victim middleware to test
	middleware := apiserv.NewLogHandlerMiddleware(
		handler,
		log.NewLogger(log.NewJSONHandler(buffer, slog.LevelInfo)),
		slog.LevelInfo,
		"http",
	)

	// act: make a request to the handler
	req := httptest.NewRequest(http.MethodGet, "http://example.com/foo?q=1", strings.NewReader("1234567890"))
	req.Header.Add("User-Agent", "test-agent")
	w := httptest.NewRecorder()
	middleware.ServeHTTP(w, req)

	// assert: the response logs should contain the request details
	t.Logf("logged: >>> %s <<<", buffer.String())
	var message map[string]interface{}
	require.NoError(t, json.Unmarshal([]byte(buffer.String()), &message), "failed to unmarshal log message")
	httpDetails, _ := message["http"].(map[string]interface{})
	assert.Equal(t, "INFO", message["level"])
	assert.Equal(t, "http", message["msg"])
	assert.Equal(t, "192.0.2.1", httpDetails["remote_addr"])
	assert.Equal(t, "HTTP/1.1", httpDetails["proto"])
	assert.Equal(t, "GET", httpDetails["method"])
	assert.Equal(t, "/foo", httpDetails["path"])
	assert.Equal(t, "q=1", httpDetails["query"])
	assert.Equal(t, "test-agent", httpDetails["user_agent"])
	assert.Equal(t, 10, mustAsInt(httpDetails["content_length"]))
	assert.Equal(t, 200, mustAsInt(httpDetails["status_code"]))
	assert.Equal(t, len(helloWorld), mustAsInt(httpDetails["response_bytes"]))
}

func TestNewLogHandlerMiddlewareErrorsShouldBeLoggedWithErrorLevel(t *testing.T) {
	t.Parallel()
	// arrange: an handler that returns "Hello World!" and sets status code to 500
	helloWorld := "<html><body>Hello World!</body></html>"
	handler := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = io.WriteString(w, helloWorld)
	})
	buffer := &strings.Builder{}
	// arrange: a victim middleware to test
	middleware := apiserv.NewLogHandlerMiddleware(
		handler,
		log.NewLogger(log.NewJSONHandler(buffer, slog.LevelInfo)),
		slog.LevelInfo,
		"http",
	)

	// act: make a request to the handler
	req := httptest.NewRequest(http.MethodGet, "http://example.com/foo?q=1", strings.NewReader("1234567890"))
	req.Header.Add("User-Agent", "test-agent")
	w := httptest.NewRecorder()
	middleware.ServeHTTP(w, req)

	// assert: the response logs should contain the request details
	t.Logf("logged: >>> %s <<<", buffer.String())
	var message map[string]interface{}
	require.NoError(t, json.Unmarshal([]byte(buffer.String()), &message), "failed to unmarshal log message")
	assert.Equal(t, "ERROR", message["level"])
}

func mustAsInt(v interface{}) int {
	vi, ok := v.(float64)
	if !ok {
		panic(fmt.Sprintf("unexpected type: %T (%v)", v, v))
	}
	return int(vi)
}
