package serviceotel

import (
	"connectrpc.com/connect"
	"connectrpc.com/otelconnect"
)

func DefaultServicesInterceptors() []connect.Interceptor {
	var interceptors []connect.Interceptor
	if interceptor, err := otelconnect.NewInterceptor(); err == nil {
		interceptors = append(interceptors, interceptor)
	}
	return interceptors
}
