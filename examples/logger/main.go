package main

import (
	"context"
	"log/slog"

	"github.com/chainguard-dev/clog"
)

func main() {
	log := clog.NewLogger(slog.Default()).With("a", "b")
	ctx := clog.WithLogger(context.Background(), log)

	// Grab logger from context and use
	clog.FromContext(ctx).With("foo", "bar").Infof("hello world")

	// Package level context loggers are also aware
	clog.ErrorContext(ctx, "asdf")
}
