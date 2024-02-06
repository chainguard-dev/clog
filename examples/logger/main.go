package main

import (
	"context"
	"flag"
	"log/slog"
	"os"

	"github.com/chainguard-dev/clog"
	"github.com/chainguard-dev/clog/slag"
)

func main() {
	var level slag.Level
	flag.Var(&level, "log-level", "log level")
	flag.Parse()

	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: &level})))

	log := clog.NewLogger(slog.Default()).With("a", "b")
	ctx := clog.WithLogger(context.Background(), log)

	// Grab logger from context and use
	clog.FromContext(ctx).With("foo", "bar").Debugf("hello debug world")
	clog.FromContext(ctx).With("info", true).Infof("hello info world")
	clog.FromContext(ctx).With("warn", 42).Warnf("hello warn world")

	// Package level context loggers are also aware
	clog.ErrorContext(ctx, "hello error world")
}
