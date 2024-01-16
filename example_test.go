package clog_test

import (
	"context"
	"log/slog"
	"os"

	"github.com/chainguard-dev/clog"
	"github.com/chainguard-dev/clog/slogtest"
)

func ExampleHandler() {
	log := slog.New(clog.NewHandler(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		// Remove time for repeatable results
		ReplaceAttr: slogtest.RemoveTime,
	})))

	ctx := context.Background()
	ctx = clog.WithValues(ctx, "foo", "bar")
	log.InfoContext(ctx, "hello world", slog.Bool("baz", true))

	// Output:
	// level=INFO msg="hello world" baz=true foo=bar
}

func ExampleLogger() {
	log := clog.NewLogger(slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		// Remove time for repeatable results
		ReplaceAttr: slogtest.RemoveTime,
	})))
	log = log.With("a", "b")
	ctx := clog.WithLogger(context.Background(), log)

	// Grab logger from context and use
	// Note: this is a formatter aware method, not an slog.Attr method.
	clog.FromContext(ctx).With("foo", "bar").Infof("hello %s", "world")

	// Package level context loggers are also aware
	clog.ErrorContext(ctx, "asdf", slog.Bool("baz", true))

	// Output:
	// level=INFO msg="hello world" a=b foo=bar
	// level=ERROR msg=asdf a=b baz=true
}
