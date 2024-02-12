package slogtest

import (
	"context"
	"io"
	"log/slog"

	"github.com/chainguard-dev/clog"
)

var (
	_ io.Writer = &logAdapter{}
)

type logAdapter struct {
	l Logger
}

func (l *logAdapter) Write(b []byte) (int, error) {
	l.l.Helper()
	l.l.Log(string(b))
	return len(b), nil
}

type Logger interface {
	Helper()
	Log(args ...any)
	Logf(format string, args ...any)
}

// TestLogger gets a logger to use in unit and end to end tests.
// This logger is configured to log at debug level.
func TestLogger(t Logger) *clog.Logger {
	return clog.New(slog.NewTextHandler(&logAdapter{l: t}, &slog.HandlerOptions{}))
}

// TestContextWithLogger returns a context with a logger to be used in tests
func TestContextWithLogger(t Logger) context.Context {
	return clog.WithLogger(context.Background(), TestLogger(t))
}

// RemoveTime removes the top-level time attribute.
// It is intended to be used as a ReplaceAttr function,
// to make example output deterministic.
//
// This is taken from slog/internal/slogtest.RemoveTime.
func RemoveTime(groups []string, a slog.Attr) slog.Attr {
	if a.Key == slog.TimeKey && len(groups) == 0 {
		return slog.Attr{}
	}
	return a
}
