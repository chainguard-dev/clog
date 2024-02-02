package main

import (
	"context"
	"flag"
	"log/slog"
	"os"

	"github.com/chainguard-dev/clog"
)

func main() {
	level := clog.LevelFlag("log-level", slog.LevelInfo, "log level")
	flag.Parse()

	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
		Level: slog.Level(level.Get()),
	})))

	log := clog.NewLogger(slog.Default()).With("a", "b")
	ctx := clog.WithLogger(context.Background(), log)

	// Grab logger from context and use
	clog.FromContext(ctx).With("foo", "bar").Infof("hello world")

	clog.FromContext(ctx).With("foo", "bar").Debugf("hello debug world")

	// Package level context loggers are also aware
	clog.ErrorContext(ctx, "asdf")
}
