package main

import (
	"context"
	"log/slog"
	"os"

	"github.com/wlynch/slogctx"
)

func init() {
	slog.SetDefault(slog.New(slogctx.NewHandler(slog.NewTextHandler(os.Stdout, nil))))
}

func main() {
	ctx := context.Background()
	ctx = slogctx.With(ctx, "foo", "bar")
	slog.InfoContext(ctx, "hello world")
}
