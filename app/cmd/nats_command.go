package cmd

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"os"
	"strings"

	"github.com/nats-io/nats.go"
	"github.com/spf13/cobra"

	"github.com/leonardinius/go-service-template/internal/apiworker"
	"github.com/leonardinius/go-service-template/internal/services"
	"github.com/leonardinius/go-service-template/internal/services/version"
)

const (
	metricsDefaultListenPort    = "8080"
	metricsDefaultListenAddress = "localhost:8080"
)

type natsCommand struct {
	c              *cobra.Command
	logLevel       string
	metricsAddress string
	//--nats--
	url      string
	user     string
	password string
	creds    string
	nkey     string
	tlscert  string
	tlskey   string
	tlsca    string
}

func CreateapiworkerCommand(context.Context) *natsCommand {
	r := natsCommand{}

	r.c = &cobra.Command{
		Use:   "nats",
		Short: "Run NATS.io worker",
		Long: "`nats` starts an NATS.io worker. Additionally exposes metrics on http://[metrics]/metrics.\n" +
			"Example:\n" +
			"\tnats --server nats://localhost:4222 --user user --password password",
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

	r.c.Flags().StringVarP(&r.metricsAddress, "metrics", "m", metricsDefaultListenAddress, "[[host]:port] listen address")
	r.c.PersistentFlags().StringVar(&r.logLevel, "log-level", "info", "log level: debug, info, warn, error")
	r.c.Flags().StringVar(&r.url, "server", nats.DefaultURL, "NATS server urls (URLs)")
	r.c.Flags().StringVar(&r.user, "user", "", "Username or Token (USERNAME)")
	r.c.Flags().StringVar(&r.password, "password", "", "Password (PASSWORD)")
	r.c.Flags().StringVar(&r.creds, "creds", "", "User credentials (FILE)")
	r.c.Flags().StringVar(&r.nkey, "nkey", "", "User NKEY (FILE)")
	r.c.Flags().StringVar(&r.tlscert, "tlscert", "", "TLS public certificate (FILE)")
	r.c.Flags().StringVar(&r.tlskey, "tlskey", "", "TLS private key (FILE)")
	r.c.Flags().StringVar(&r.tlsca, "tlsca", "", "TLS certificate authority chain (FILE)")
	return &r
}

func (r *natsCommand) Command() *cobra.Command {
	return r.c
}

func (r *natsCommand) execute(ctx context.Context) error {
	return r.runServe(ctx, r.metricsAddress)
}

func (r *natsCommand) runServe(ctx context.Context, metricsAddress string) error {
	// check if address is host/ip:port
	_, _, err := net.SplitHostPort(metricsAddress)
	if err != nil && strings.Contains(err.Error(), "missing port in address") {
		// if there is no port, append default port
		metricsAddress = net.JoinHostPort(metricsAddress, metricsDefaultListenPort)
	}

	slog.LogAttrs(ctx, slog.LevelInfo, "starting nats worker",
		slog.String("version", version.FullVersion),
		slog.String("metrics", metricsAddress),
		slog.String("server", r.url),
		slog.String("user", r.user),
		slog.String("password", strings.Repeat("*", len(r.password))),
		slog.String("creds", r.creds),
		slog.String("nkey", r.nkey),
		slog.String("tlscert", r.tlscert),
		slog.String("tlskey", r.tlskey),
		slog.String("tlsca", r.tlsca))

	url := r.url
	if r.url == "" {
		url = nats.DefaultURL
	}
	options := r.natsOptions()

	wrk, err := apiworker.NewWorker(ctx, url, services.AllRoutes, options...)
	if err != nil {
		return err
	}

	errCh := make(chan error, 1)
	go func() {
		errCh <- wrk.ListenAndServe(ctx)
	}()

	select {
	case err := <-errCh:
		return err
	case <-ctx.Done():
		slog.LogAttrs(ctx, slog.LevelInfo, "shutting down nats worker...", slog.String("server", url))
		return wrk.Shutdown(context.WithoutCancel(ctx))
	}
}

func (r *natsCommand) natsOptions() []apiworker.Option {
	var options []apiworker.Option
	if r.metricsAddress != "" {
		options = append(options, apiworker.WithMetricsAddress(r.metricsAddress))
	}
	if r.url != "" {
		options = append(options, apiworker.WithURL(r.url))
	}
	if r.user != "" {
		options = append(options, apiworker.WithUser(r.user))
	}
	if r.password != "" {
		options = append(options, apiworker.WithPassword(r.password))
	}
	if r.creds != "" {
		creds := mustExpandPath(r.creds)
		options = append(options, apiworker.WithCreds(creds))
	}
	if r.nkey != "" {
		nkey := mustExpandPath(r.nkey)
		options = append(options, apiworker.WithNKey(nkey))
	}
	if r.tlscert != "" {
		tlscert := mustExpandPath(r.tlscert)
		options = append(options, apiworker.WithTLSCert(tlscert))
	}
	if r.tlskey != "" {
		tlskey := mustExpandPath(r.tlskey)
		options = append(options, apiworker.WithTLSKey(tlskey))
	}
	if r.tlsca != "" {
		tlsca := mustExpandPath(r.tlsca)
		options = append(options, apiworker.WithTLSCA(tlsca))
	}
	return options
}

func mustExpandPath(path string) string {
	path, err := expandPath(path)
	if err != nil {
		panic(fmt.Errorf("failed to expand path %q: %w", path, err))
	}
	return path
}

func expandPath(path string) (string, error) {
	if strings.HasPrefix(path, "~") {
		home, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}
		path = home + path[1:]
	}

	path = os.ExpandEnv(path)
	return path, nil
}
