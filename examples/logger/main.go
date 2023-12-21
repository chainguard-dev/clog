package main

import (
	"context"
	"log/slog"

	"github.com/chainguard-dev/slogctx"
)

func main() {
	log := slogctx.NewLogger(slog.Default()).With("a", "b")
	ctx := slogctx.WithLogger(context.Background(), log)

	// Grab logger from context and use
	slogctx.FromContext(ctx).With("foo", "bar").Infof("hello world")

	// Package level context loggers are also aware
	slogctx.ErrorContext(ctx, "asdf")
}
