package slogtest

import (
	"context"
	"io"
	"log/slog"
	"strings"

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
	l.l.Log(strings.TrimSuffix(string(b), "\n"))
	return len(b), nil
}

type Logger interface {
	Log(args ...any)
	Logf(format string, args ...any)
	Helper()
}

// TestLogger gets a logger to use in unit and end to end tests.
// This logger is configured to log at debug level.
func TestLogger(t Logger) *clog.Logger {
	return clog.New(slog.NewTextHandler(&logAdapter{l: t}, &slog.HandlerOptions{AddSource: true}))
}

// TestLoggerWithOptions gets a logger to use in unit and end to end tests.
func TestLoggerWithOptions(t Logger, opts *slog.HandlerOptions) *clog.Logger {
	return clog.New(slog.NewTextHandler(&logAdapter{l: t}, opts))
}

// Context returns a context with a logger to be used in tests.
// This is equivalent to TestContextWithLogger.
func Context(t Logger) context.Context {
	return TestContextWithLogger(t)
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
