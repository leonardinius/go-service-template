package servehttpe2e_test

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"testing"

	"connectrpc.com/connect"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/leonardinius/go-service-template/app/cmd"
	"github.com/leonardinius/go-service-template/internal/apigen/version/v1/versionv1connect"
	"github.com/leonardinius/go-service-template/internal/insights"
	"github.com/leonardinius/go-service-template/internal/services/version"
	"github.com/leonardinius/go-service-template/teste2e/internal/testbind"
	"github.com/leonardinius/go-service-template/teste2e/internal/testhttp"

	versionv1 "github.com/leonardinius/go-service-template/internal/apigen/version/v1"
)

func TestServeHTTPVersionInfoConnectHTTP(t *testing.T) {
	t.Parallel()
	runTest(t, func(ctx context.Context, port int) {
		resp := testhttp.MustPost(ctx, t,
			endpointURL("http://localhost:{{port}}", port, versionv1connect.VersionServiceGetVersionProcedure),
			"application/json",
			strings.NewReader("{}"))
		require.Equal(t, 200, resp.StatusCode, "Expected 200 OK, got %s", resp.Status)
		assert.Contains(t, resp.Header.Get("Content-Type"), "application/json")
		var contents map[string]interface{}
		testhttp.MustReadFullyJSON(t, resp, &contents)
		versionResponse, ok := contents["version"].(map[string]interface{})
		require.True(t, ok, "Expected 'version' field in the response")
		require.Equal(t, version.FullVersion, versionResponse["fullVersion"])
		_ = resp.Body.Close()
	})
}

func TestServeHTTPVersionInfoGrpc(t *testing.T) {
	t.Parallel()
	runTest(t, func(ctx context.Context, port int) {
		client := versionv1connect.NewVersionServiceClient(&http.Client{},
			endpointURL("http://localhost:{{port}}", port),
			connect.WithGRPC(),
		)

		resp, err := client.GetVersion(ctx, connect.NewRequest(&versionv1.GetVersionRequest{}))

		require.NoError(t, err)
		versionResponse := resp.Msg.GetVersion()
		require.Equal(t, version.FullVersion, versionResponse.GetFullVersion())
	})
}

func TestServeHTTPMetrics(t *testing.T) {
	t.Parallel()
	runTest(t, func(ctx context.Context, port int) {
		resp := testhttp.MustPost(ctx, t,
			endpointURL("http://localhost:{{port}}", port, versionv1connect.VersionServiceGetVersionProcedure),
			"application/json",
			strings.NewReader("{}"))
		_ = resp.Body.Close()

		resp = testhttp.MustGET(ctx, t, endpointURL("http://localhost:{{port}}/metrics", port))
		require.Equal(t, 200, resp.StatusCode, "Expected 200 OK, got %s", resp.Status)
		contents := testhttp.MustReadFullyString(t, resp)
		assert.Contains(t, contents, "# HELP go_info Information about the Go environment")
		assert.Contains(t, contents, "# TYPE go_info gauge")
		assert.Contains(t, contents, "go_info{")
		assert.Contains(t, contents, versionv1connect.VersionServiceGetVersionProcedure)
		_ = resp.Body.Close()
	})
}

func TestOtelHasXTraceHeadersConnectHttp(t *testing.T) {
	t.Parallel()
	runTest(t, func(ctx context.Context, port int) {
		resp := testhttp.MustGET(ctx, t, endpointURL("http://localhost:{{port}}", port))
		xTraceID := resp.Header.Get("X-Trace-Id")
		assert.NotEmpty(t, xTraceID)
		assert.Len(t, xTraceID, 32, "Expected 32 characters, got %d", len(xTraceID))

		traceParent := resp.Header.Get("Traceparent")
		assert.NotEmpty(t, traceParent)
		assert.Len(t, traceParent, 55, "Expected 55 characters, got %d", len(xTraceID))
		_ = resp.Body.Close()
	})
}

func TestOtelHasXTraceHeadersGrpc(t *testing.T) {
	t.Parallel()
	runTest(t, func(ctx context.Context, port int) {
		client := versionv1connect.NewVersionServiceClient(&http.Client{},
			endpointURL("http://localhost:{{port}}", port),
			connect.WithGRPC(),
		)

		resp, err := client.GetVersion(ctx, connect.NewRequest(&versionv1.GetVersionRequest{}))
		require.NoError(t, err)
		xTraceID := resp.Header().Get("X-Trace-Id")
		assert.NotEmpty(t, xTraceID)
		assert.Len(t, xTraceID, 32, "Expected 32 characters, got %d", len(xTraceID))

		traceParent := resp.Header().Get("Traceparent")
		assert.NotEmpty(t, traceParent)
		assert.Len(t, traceParent, 55, "Expected 55 characters, got %d", len(xTraceID))
	})
}

func runTest(t *testing.T, test func(ctx context.Context, port int)) {
	t.Helper()

	port := testbind.DynamicPort()

	ctx := context.WithoutCancel(rootTestCtx)
	ctx = insights.ContextWithRegistry(ctx, insights.NewMetricsRegistry())
	ctx, stopMain := context.WithCancel(ctx)
	defer stopMain()

	address := fmt.Sprintf("localhost:%d", port)
	errCh := make(chan error, 1)
	serveCommand := cmd.CreateHTTPServeCommand(ctx)
	serveCommand.Command().SetArgs([]string{
		"--address=" + address,
	})
	go func() {
		errCh <- serveCommand.Command().ExecuteContext(ctx)
	}()

	testbind.MustWaitForPortListenUp(ctx, t, port)
	test(ctx, port)
	stopMain()
	// Full shutdown of the server may take 3 more seconds.
	// Uncomment this line if there is a need to check the return error.
	// <-errCh
	//
	// Instead, we can just wait for the port to be taken down.
	// It is documented the ListenAndServe function will return
	// immediately after the server is shutdown (on signal Done).
	testbind.MustWaitForPortListenDown(ctx, t, port)
}

func endpointURL(url string, port int, parts ...string) string {
	base := strings.ReplaceAll(url, "{{port}}", strconv.Itoa(port))
	return base + strings.Join(parts, "/")
}

var rootTestCtx = context.Background()

func TestMain(m *testing.M) {
	cmd.MustSetupLogger(rootTestCtx, "info")
	otelShutdown, err := insights.SetupOtelSDK(rootTestCtx)
	if err != nil {
		panic(err)
	}
	m.Run()
	err = otelShutdown(rootTestCtx)
	if err != nil {
		panic(err)
	}
}
