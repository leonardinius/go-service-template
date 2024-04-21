package testbind

import (
	"net"
	"sync"
)

var acquiredPorts = new(sync.Map)

// DynamicPort supplies random free net ports to use.
func DynamicPort() int {
	port := mustBindDynamicPort()
	for {
		if _, loaded := acquiredPorts.LoadOrStore(port, port); !loaded {
			break
		}
		port = mustBindDynamicPort()
	}
	return port
}

func mustBindDynamicPort() int {
	listener, err := net.Listen("tcp", LocalhostBindHost.AddrWithPort(0))
	if err != nil {
		panic(err)
	}
	defer func() {
		err := listener.Close()
		if err != nil {
			panic(err)
		}
	}()

	addr, ok := listener.Addr().(*net.TCPAddr)
	if !ok {
		panic("listener.Addr() is not *net.TCPAddr")
	}

	return addr.Port
}
