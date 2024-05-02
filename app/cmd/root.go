package cmd

import (
	"context"
	"errors"

	"github.com/spf13/cobra"

	"github.com/leonardinius/go-service-template/internal/insights"
	"github.com/leonardinius/go-service-template/internal/services/version"
)

type rootCommand struct {
	c *cobra.Command
}

func CreateRootCommand(context.Context) *rootCommand {
	r := rootCommand{}
	r.c = &cobra.Command{
		Use:     version.ServiceName,
		Version: version.FullVersion,

		// service server
		Short: "{HTTP,gRPC,NATS.io} server with Opentelemetry and Prometheus metrics, etc.",

		CompletionOptions: cobra.CompletionOptions{
			HiddenDefaultCmd: true,
		},
	}

	return &r
}

func (r *rootCommand) Command() *cobra.Command {
	return r.c
}

func Execute(ctx context.Context, args []string) error {
	otelShutdown, err := insights.SetupOtelSDK(ctx)
	if err != nil {
		return err
	}

	// Handle shutdown properly so nothing leaks.
	otelHandleShutdown := func(e error) error {
		return errors.Join(e, otelShutdown(ctx))
	}

	rootCmd := CreateRootCommand(ctx).Command()
	rootCmd.AddCommand(CreateHTTPServeCommand(ctx).Command())
	rootCmd.AddCommand(CreateApiworkerCommand(ctx).Command())

	rootCmd.SetArgs(args)

	return otelHandleShutdown(rootCmd.ExecuteContext(ctx))
}
