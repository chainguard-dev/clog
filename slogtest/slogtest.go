// Package slogtest provides utilities for emitting test logs using clog.
//
//	func TestExample(t *testing.T) {
//		ctx := slogtest.Context(t)
//		clog.FromContext(ctx).With("foo", "bar").Info("hello world")
//	}
//
// This produces the following test output:
//
//	=== RUN   TestExample
//		slogtest.go:24: level=INFO source=/path/to/example_test.go:13 msg="hello world" foo=bar
//
// This package is intended to be used in tests only.
package slogtest

import (
	"context"
	"io"
	"log/slog"
	"strings"
	"testing"

	"github.com/chainguard-dev/clog"
)

var _ io.Writer = &logAdapter{}

type logAdapter struct{ l Logger }

func (l *logAdapter) Write(b []byte) (int, error) {
	l.l.Log(strings.TrimSuffix(string(b), "\n"))
	return len(b), nil
}

var _ Logger = (*testing.T)(nil)
var _ Logger = (*testing.B)(nil)
var _ Logger = (*testing.F)(nil)

type Logger interface{ Log(args ...any) }

// TestLogger gets a logger to use in unit and end to end tests.
// This logger is configured to log at debug level.
func TestLogger(t Logger) *clog.Logger {
	return clog.New(slog.NewTextHandler(&logAdapter{l: t}, &slog.HandlerOptions{
		Level:       slog.LevelDebug,
		AddSource:   true,
		ReplaceAttr: RemoveTime,
	}))
}

// TestLoggerWithOptions gets a logger to use in unit and end to end tests.
func TestLoggerWithOptions(t Logger, opts *slog.HandlerOptions) *clog.Logger {
	return clog.New(slog.NewTextHandler(&logAdapter{l: t}, opts))
}

// Context returns a context with a logger to be used in tests.
func Context(t Logger) context.Context {
	return clog.WithLogger(context.Background(), TestLogger(t))
}

// TestContextWithLogger returns a context with a logger to be used in tests
//
// Deprecated: Use Context instead.
func TestContextWithLogger(t Logger) context.Context { return Context(t) }

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
