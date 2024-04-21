package cmd

import (
	"context"
	"log/slog"
	"net"
	"strings"

	"github.com/spf13/cobra"

	"github.com/leonardinius/go-service-template/internal/apiserv"
	"github.com/leonardinius/go-service-template/internal/services"
	"github.com/leonardinius/go-service-template/internal/services/version"
)

const (
	httpDefaultListenPort    = "8080"
	httpDefaultListenAddress = "localhost:8080"
)

type httpCommand struct {
	c             *cobra.Command
	logLevel      string
	listenAddress string
}

func CreateHTTPServeCommand(context.Context) *httpCommand {
	r := httpCommand{}

	r.c = &cobra.Command{
		Use:   "http",
		Short: "Run {HTTP, gRPC} server",
		Long: "`http` starts an HTTP server on the specified address and port. Default is " + httpDefaultListenAddress + "." +
			"\n" +
			"The server is HTTP & gRPC-compatible (see https://connectrpc.com for more details).",
		//nolint:contextcheck // cobra interface
		PreRun: func(cmd *cobra.Command, args []string) {
			// Do not print usage on error, eg when port is already in use.
			cmd.SilenceUsage = true
			MustSetupLogger(cmd.Context(), r.logLevel)
		},

		//nolint:contextcheck // cobra interface
		RunE: func(cmd *cobra.Command, args []string) error {
			return r.execute(cmd.Context())
		},
		Args: cobra.NoArgs,
	}

	r.c.Flags().StringVarP(&r.listenAddress, "address", "a", httpDefaultListenAddress, "[[host]:port] listen address")
	r.c.PersistentFlags().StringVar(&r.logLevel, "log-level", "info", "log level: debug, info, warn, error")
	return &r
}

func (r *httpCommand) Command() *cobra.Command {
	return r.c
}

func (r *httpCommand) execute(ctx context.Context) error {
	return r.runServe(ctx, r.listenAddress)
}

func (r *httpCommand) runServe(ctx context.Context, address string) error {
	// check if address is host/ip:port
	_, _, err := net.SplitHostPort(address)
	if err != nil && strings.Contains(err.Error(), "missing port in address") {
		// if there is no port, append default port
		address = net.JoinHostPort(address, httpDefaultListenPort)
	}

	go func() {
		<-ctx.Done()
		slog.LogAttrs(ctx, slog.LevelInfo, "signal received, shutting down...", slog.String("address", address))
	}()

	slog.LogAttrs(ctx, slog.LevelInfo, "starting http",
		slog.String("version", version.FullVersion),
		slog.String("address", address))

	srv, err := apiserv.NewDefaultServer(ctx, address, services.AllRoutes...)
	if err != nil {
		return err
	}

	errCh := make(chan error, 1)
	go func() {
		errCh <- apiserv.ListenAndServe(ctx, srv)
	}()

	select {
	case err := <-errCh:
		return err
	case <-ctx.Done():
		slog.LogAttrs(ctx, slog.LevelInfo, "shutting down http server")
		return srv.Shutdown(context.WithoutCancel(ctx))
	}
}
