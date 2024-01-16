package main

import (
	"context"
	"log/slog"
	"os"

	"github.com/chainguard-dev/clog"
)

func init() {
	slog.SetDefault(slog.New(clog.NewHandler(slog.NewTextHandler(os.Stdout, nil))))
}

func main() {
	ctx := context.Background()
	ctx = clog.WithValues(ctx, "foo", "bar")

	// Use slog package directly
	slog.InfoContext(ctx, "hello world", slog.Bool("baz", true))

	// glog / zap style (note: can't pass additional attributes)
	clog.Errorf("hello %s", "world")
}
