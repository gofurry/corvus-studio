package main

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/signal"
	"syscall"

	"github.com/gofurry/corvus-studio/apps/core/internal/application"
	"github.com/gofurry/corvus-studio/apps/core/internal/interfaces/cli"
)

func run(ctx context.Context, args []string, stdout, stderr io.Writer) error {
	command := cli.NewRootCommand(cli.Dependencies{
		Run: application.Run,
	})
	command.SetArgs(args)
	command.SetOut(stdout)
	command.SetErr(stderr)
	return command.ExecuteContext(ctx)
}

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	if err := run(ctx, os.Args[1:], os.Stdout, os.Stderr); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, "corvus:", err)
		os.Exit(1)
	}
}
