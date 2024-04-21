package testbind

import (
	"fmt"
	"time"
)

var (
	PortPollInterval = 10 * time.Millisecond
	PortPollTimeout  = 3 * time.Second
)

type bindHost string

const LocalhostBindHost = bindHost("localhost")

var _ fmt.Stringer = LocalhostBindHost

func (h bindHost) String() string {
	return string(h)
}

func (h bindHost) AddrWithPort(port int) string {
	return fmt.Sprintf("%s:%d", h, port)
}
