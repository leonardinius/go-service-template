package apiserv_test

import (
	"encoding/json"
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

func TestNewRecoveryHandlerMiddlewareSmokeTest(t *testing.T) {
	t.Parallel()
	// arrange: an handler that panics
	var handler http.Handler = http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		panic("hello, panic!")
	})
	buffer := &strings.Builder{}
	// arrange: a victim middleware to test
	middleware := apiserv.NewRecoveryHandlerMiddleware(
		handler,
		log.NewLogger(log.NewJSONHandler(buffer, slog.LevelInfo)),
	)

	// act: make a request to the handler
	req := httptest.NewRequest(http.MethodGet, "http://example.com/foo?q=1", strings.NewReader("1234567890"))
	req.Header.Add("User-Agent", "test-agent")
	w := httptest.NewRecorder()
	middleware.ServeHTTP(w, req)

	// assert: the response should be 500
	assert.Equal(t, 500, w.Code)
	assert.Equal(t, "Internal Server Error\n", w.Body.String())

	// assert: the log message should contain the panic details
	var message map[string]interface{}
	require.NoError(t, json.Unmarshal([]byte(buffer.String()), &message), "failed to unmarshal log message")
	assert.Equal(t, "ERROR", message["level"])
	assert.Equal(t, "recovered from panic", message["msg"])
	assert.Equal(t, "panic caught: hello, panic!", message["error"])
	assert.Contains(t, message["stack"], t.Name())
}
