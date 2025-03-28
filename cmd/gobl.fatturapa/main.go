package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

// build data provided by goreleaser and mage setup
var (
	name    = "gobl.fatturapa"
	version = "dev"
	date    = ""
)

func main() {
	if err := run(); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run() error {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	return root().cmd().ExecuteContext(ctx)
}

func inputFilename(args []string) string {
	if len(args) > 0 && args[0] != "-" {
		return args[0]
	}
	return ""
}

func outputFilename(args []string) string {
	if len(args) >= 2 && args[1] != "-" {
		return args[1]
	}
	return ""
}
