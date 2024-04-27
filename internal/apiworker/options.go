package apiworker

import (
	"context"
	"log/slog"

	natsio "github.com/nats-io/nats.go"
)

type natsOptions struct {
	metricsAddress string
	url            string
	user           string
	password       string
	creds          string
	nkey           string
	tlscert        string
	tlskey         string
	tlsca          string
	context        string
}

type Option interface {
	apply(options *natsOptions) error
}

type funcOption func(*natsOptions) error

func (fo funcOption) apply(o *natsOptions) error {
	return fo(o)
}

func WithMetricsAddress(address string) Option {
	return funcOption(func(o *natsOptions) error {
		o.metricsAddress = address
		return nil
	})
}

func WithURL(url string) Option {
	return funcOption(func(o *natsOptions) error {
		o.url = url
		return nil
	})
}

func WithUser(user string) Option {
	return funcOption(func(o *natsOptions) error {
		o.user = user
		return nil
	})
}

func WithPassword(password string) Option {
	return funcOption(func(o *natsOptions) error {
		o.password = password
		return nil
	})
}

func WithCreds(creds string) Option {
	return funcOption(func(o *natsOptions) error {
		o.creds = creds
		return nil
	})
}

func WithNKey(nkey string) Option {
	return funcOption(func(o *natsOptions) error {
		o.nkey = nkey
		return nil
	})
}

func WithTLSCert(tlscert string) Option {
	return funcOption(func(o *natsOptions) error {
		o.tlscert = tlscert
		return nil
	})
}

func WithTLSKey(tlskey string) Option {
	return funcOption(func(o *natsOptions) error {
		o.tlskey = tlskey
		return nil
	})
}

func WithTLSCA(tlsca string) Option {
	return funcOption(func(o *natsOptions) error {
		o.tlsca = tlsca
		return nil
	})
}

func WithContextName(contextName string) Option {
	return funcOption(func(o *natsOptions) error {
		o.context = contextName
		return nil
	})
}

func newNatsOptions(opts ...Option) (*natsOptions, error) {
	options := &natsOptions{}
	for _, opt := range opts {
		err := opt.apply(options)
		if err != nil {
			return nil, err
		}
	}
	return options, nil
}

func buildNatsIOOptions(ctx context.Context, opts ...Option) (*natsOptions, []natsio.Option, error) {
	config, err := newNatsOptions(opts...)
	if err != nil {
		return nil, nil, err
	}
	var natsioOptions []natsio.Option = nil
	if config.user != "" {
		natsioOptions = append(natsioOptions, natsio.UserInfo(config.user, config.password))
	}
	if config.creds != "" {
		natsioOptions = append(natsioOptions, natsio.UserCredentials(config.creds))
	}
	if config.nkey != "" {
		if opt, err := natsio.NkeyOptionFromSeed(config.nkey); err == nil {
			natsioOptions = append(natsioOptions, opt)
		} else {
			return nil, nil, err
		}
	}
	if config.tlscert != "" || config.tlskey != "" {
		natsioOptions = append(natsioOptions, natsio.ClientCert(config.tlscert, config.tlskey))
	}
	if config.tlsca != "" {
		natsioOptions = append(natsioOptions, natsio.RootCAs(config.tlsca))
	}

	natsConCb := func(name string) natsio.ConnHandler {
		return func(nc *natsio.Conn) {
			slog.LogAttrs(ctx, slog.LevelInfo, "natsio connection", slog.String("event", name))
		}
	}
	disconnectErrCb := func(nc *natsio.Conn, err error) {
		slog.LogAttrs(ctx, slog.LevelError, "natsio disconnected", slog.Any("error", err))
	}
	errorCb := func(nc *natsio.Conn, sub *natsio.Subscription, err error) {
		slog.LogAttrs(ctx, slog.LevelError, "natsio error", slog.Any("error", err))
	}
	natsioOptions = append(natsioOptions,
		natsio.ConnectHandler(natsConCb("connect")),
		// natsio.ClosedHandler(natsConCb("closed")),
		natsio.DisconnectHandler(natsConCb("disconnect")),
		natsio.ReconnectHandler(natsConCb("reconnect")),
		natsio.DiscoveredServersHandler(natsConCb("discovered_servers")),
		natsio.DisconnectErrHandler(disconnectErrCb),
		natsio.ErrorHandler(errorCb),
	)

	return config, natsioOptions, nil
}
