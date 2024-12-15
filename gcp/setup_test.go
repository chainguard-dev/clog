package gcp

import (
	"context"
	"log/slog"
	"testing"
)

func TestHandler(t *testing.T) {
	ctx := context.Background()
	l := slog.New(NewHandler(slog.LevelInfo))
	l.With("level", "INFO").Log(ctx, slog.LevelInfo, "hello world") // okay
	l.With("level", "INFO").Log(ctx, slog.LevelWarn, "hello world") // weird, but okay (info)

	// These should not panic.
	l.With("level", nil).Log(ctx, slog.LevelInfo, "hello world")
	l.With("level", 123).Log(ctx, slog.LevelInfo, "hello world")
	l.With("level", map[string]string{}).Log(ctx, slog.LevelInfo, "hello world")
}
