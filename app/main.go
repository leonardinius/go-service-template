package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/joho/godotenv"

	"github.com/leonardinius/go-service-template/app/cmd"
)

func main() {
	os.Exit(Main(os.Args[1:]))
}

func Main(args []string) int {
	ctx := context.Background()
	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt, syscall.SIGTERM)
	defer cancel()

	err := godotenv.Load()
	if err != nil {
		slog.ErrorContext(ctx, "Error loading .env file")
		return 1
	}

	if err := cmd.Execute(ctx, args); err != nil {
		fmt.Println(err)
		return 1
	}

	return 0
}
