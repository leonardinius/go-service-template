package testnats

import (
	"context"
	"net"
	"strconv"
	"strings"
	"testing"

	"github.com/nats-io/nats.go"
	"github.com/stretchr/testify/require"
)

func MustConnect(t *testing.T, ctx context.Context, port int) *nats.Conn {
	t.Helper()

	address := "nats://" + net.JoinHostPort(Host, strconv.Itoa(port))
	nc, err := nats.Connect(address)
	require.NoErrorf(t, err, "Failed to connect to NATS server: %v", err)
	t.Cleanup(nc.Close)
	return nc
}

func Request(t *testing.T, ctx context.Context, nc *nats.Conn, path string, payload []byte) (*nats.Msg, error) {
	t.Helper()

	subj := pathToSubject(path)
	return nc.RequestWithContext(ctx, subj, payload)
}

func MustRequest(t *testing.T, ctx context.Context, nc *nats.Conn, path string, payload []byte) *nats.Msg {
	t.Helper()

	msg, err := Request(t, ctx, nc, path, payload)
	require.NoErrorf(t, err, "Failed to request %q: %v", path, err)
	return msg
}

func pathToSubject(path string) string {
	subj := strings.TrimLeft(path, "/")
	subj = strings.ReplaceAll(subj, "/", ".")
	return subj
}
