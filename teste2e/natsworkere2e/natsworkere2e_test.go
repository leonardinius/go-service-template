package natsworkere2e_test

import (
	"context"
	"net"
	"strconv"
	"strings"
	"testing"

	"github.com/nats-io/nats.go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/encoding/protojson"

	"github.com/leonardinius/go-service-template/app/cmd"
	"github.com/leonardinius/go-service-template/internal/apigen/version/v1/versionv1connect"
	"github.com/leonardinius/go-service-template/internal/insights"
	"github.com/leonardinius/go-service-template/internal/services/version"
	"github.com/leonardinius/go-service-template/teste2e/internal/testbind"
	"github.com/leonardinius/go-service-template/teste2e/internal/testhttp"
	"github.com/leonardinius/go-service-template/teste2e/internal/testnats"

	versionv1 "github.com/leonardinius/go-service-template/internal/apigen/version/v1"
)

const (
	Host = testnats.Host
)

func TestVersionRequestReplyNATS(t *testing.T) {
	t.Parallel()
	runTest(t, func(ctx context.Context, port, _ int) {
		nc := testnats.MustConnect(t, ctx, port)
		reply := testnats.MustRequest(t, ctx, nc,
			versionv1connect.VersionServiceGetVersionProcedure,
			[]byte("{}"))
		resp := versionv1.GetVersionResponse{}
		err := protojson.Unmarshal(reply.Data, &resp)
		require.NoErrorf(t, err, "Failed to unmarshal response: %v", err)
		assert.Equal(t, version.FullVersion, resp.GetVersion().GetFullVersion())
	})
}

func TestVersionRequest404ReplyErrorNATS(t *testing.T) {
	t.Parallel()
	runTest(t, func(ctx context.Context, port, _ int) {
		nc := testnats.MustConnect(t, ctx, port)
		_, err := testnats.Request(t, ctx, nc,
			"404",
			[]byte("{}"))
		require.Error(t, err)
		require.ErrorIs(t, err, nats.ErrNoResponders)
	})
}

func TestOtelHasXTraceHeadersNATSReplyMessage(t *testing.T) {
	t.Parallel()
	runTest(t, func(ctx context.Context, port, _ int) {
		nc := testnats.MustConnect(t, ctx, port)
		reply := testnats.MustRequest(t, ctx, nc,
			versionv1connect.VersionServiceGetVersionProcedure,
			[]byte("{}"))
		xTraceID := reply.Header.Get("X-Trace-Id")
		assert.NotEmpty(t, xTraceID)
		assert.Len(t, xTraceID, 32, "Expected 32 characters, got %d", len(xTraceID))

		traceParent := reply.Header.Get("Traceparent")
		assert.NotEmpty(t, traceParent)
		assert.Len(t, traceParent, 55, "Expected 55 characters, got %d", len(xTraceID))
	})
}

func TestServeNATSMetrics(t *testing.T) {
	t.Parallel()
	runTest(t, func(ctx context.Context, port, metricsPort int) {
		nc := testnats.MustConnect(t, ctx, port)
		_ = testnats.MustRequest(t, ctx, nc,
			versionv1connect.VersionServiceGetVersionProcedure,
			[]byte("{}"))

		resp := testhttp.MustGET(ctx, t, endpointURL("http://localhost:{{port}}/metrics", metricsPort))
		require.Equal(t, 200, resp.StatusCode, "Expected 200 OK, got %s", resp.Status)
		contents := testhttp.MustReadFullyString(t, resp)
		assert.Contains(t, contents, "# HELP go_info Information about the Go environment")
		assert.Contains(t, contents, "# TYPE go_info gauge")
		assert.Contains(t, contents, "go_info{")
		assert.Contains(t, contents, versionv1connect.VersionServiceGetVersionProcedure)
		_ = resp.Body.Close()
	})
}

func runTest(t *testing.T, test func(ctx context.Context, natsPort, metricsPort int)) {
	t.Helper()

	natsServerPort := testbind.DynamicPort()
	metricsPort := testbind.DynamicPort()

	ctx := context.WithoutCancel(rootTestCtx)
	ctx = insights.ContextWithRegistry(ctx, insights.NewMetricsRegistry())
	ctx, stopMain := context.WithCancel(ctx)
	defer stopMain()

	ns := testnats.MustRunNatsServer(t, ctx, natsServerPort)
	defer ns.Shutdown()

	address := "nats://" + net.JoinHostPort(Host, strconv.Itoa(natsServerPort))
	metricsAddress := net.JoinHostPort(Host, strconv.Itoa(metricsPort))
	errCh := make(chan error, 1)
	serveCommand := cmd.CreateNATSWorkerCommand(ctx)
	serveCommand.Command().SetArgs([]string{
		"--server=" + address,
		"--metrics=" + metricsAddress,
	})
	go func() {
		errCh <- serveCommand.Command().ExecuteContext(ctx)
	}()
	testbind.MustWaitForPortListenUp(ctx, t, metricsPort)
	test(ctx, natsServerPort, metricsPort)
	stopMain()
	// Full shutdown of the server may take 3 more seconds.
	// Uncomment this line if there is a need to check the return error.
	// <-errCh
	//
	// Instead, we can just wait for the port to be taken down.
	// It is documented the ListenAndServe function will return
	// immediately after the server is shutdown (on signal Done).
	testbind.MustWaitForPortListenDown(ctx, t, natsServerPort)
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
