package testnats

import (
	"context"
	"testing"
	"time"

	"github.com/nats-io/nats-server/v2/server"
	"github.com/stretchr/testify/require"

	"github.com/leonardinius/go-service-template/teste2e/internal/testbind"
)

const (
	Host                     = "localhost"
	connectionTimeoutSeconds = 1
)

func MustRunNatsServer(t *testing.T, ctx context.Context, port int) *server.Server {
	t.Helper()

	opts := &server.Options{
		Port: port,
		Host: Host,
	}
	ns, err := server.NewServer(opts)
	require.NoErrorf(t, err, "Failed to create NATS server: %v", err)
	go ns.Start()

	ok := ns.ReadyForConnections(connectionTimeoutSeconds * time.Second)
	require.True(t, ok, "NATS server is not ready for connections")
	testbind.MustWaitForPortListenUp(ctx, t, port)
	return ns
}
