package testbind

import (
	"context"
	"net"
	"testing"
	"time"
)

func MustWaitForPortListenUp(ctx context.Context, t *testing.T, port int) (conn net.Conn) {
	t.Helper()
	return MustWaitForAddrPortUp(ctx, t, LocalhostBindHost.AddrWithPort(port))
}

func MustWaitForAddrPortUp(ctx context.Context, t *testing.T, addr string) (conn net.Conn) {
	t.Helper()
	poll := time.NewTicker(PortPollInterval)
	defer poll.Stop()

	timeout := time.NewTimer(PortPollTimeout)
	defer timeout.Stop()

	for {
		select {
		case <-timeout.C:
			t.Fatalf("%s port open timeout", addr)

		case <-ctx.Done():
			t.Fatalf("%s port open error, parent context is done: %v", addr, ctx.Err())

		case <-poll.C:
			var d net.Dialer
			conn, _ = checkAddr(ctx, &d, addr)
			if conn != nil {
				return conn
			}
		}
	}
}

func MustWaitForPortListenDown(ctx context.Context, t *testing.T, port int) {
	t.Helper()
	MustWaitForAddrPortDown(ctx, t, LocalhostBindHost.AddrWithPort(port))
}

func MustWaitForAddrPortDown(ctx context.Context, t *testing.T, addr string) {
	t.Helper()
	poll := time.NewTicker(PortPollInterval)
	defer poll.Stop()

	timeout := time.NewTimer(PortPollTimeout)
	defer timeout.Stop()

	awaitFailedTimes := 2
	for {
		select {
		case <-timeout.C:
			t.Fatalf("%s port close timeout", addr)

		case <-ctx.Done():
			return

		case <-poll.C:
			var d net.Dialer
			_, err := checkAddr(ctx, &d, addr)
			if err != nil {
				awaitFailedTimes--
			}
			if awaitFailedTimes == 0 {
				return
			}
		}
	}
}

func checkAddr(ctx context.Context, d *net.Dialer, addr string) (conn net.Conn, err error) {
	limitCtx, limitCancelFn := context.WithTimeout(ctx, 50*time.Millisecond)
	defer limitCancelFn()
	conn, err = d.DialContext(limitCtx, "tcp", addr)
	return conn, err
}
