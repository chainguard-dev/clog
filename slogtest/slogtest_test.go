package slogtest_test

import (
	"log/slog"
	"testing"

	"github.com/chainguard-dev/clog"
	"github.com/chainguard-dev/clog/slogtest"
)

func TestSlogTest(t *testing.T) {
	ctx := slogtest.Context(t)

	clog.FromContext(ctx).With("foo", "bar").Infof("hello world")
	clog.FromContext(ctx).With("bar", "baz").Infof("me again")
	clog.FromContext(ctx).With("baz", true).Infof("okay last one")

	clog.FromContext(ctx).Debug("hello debug")
	clog.FromContext(ctx).Info("hello info")
	clog.FromContext(ctx).Warn("hello warn")
	clog.FromContext(ctx).Error("hello error")

	fn(ctx)
}

// TestSlogTestTContext tests the use of t.Context() in Go 1.24+.
func TestSlogTestTContext(t *testing.T) {
	ctx := t.Context()
	slog.SetDefault(slog.New(slogtest.TestLogger(t).Handler()))

	clog.FromContext(ctx).With("foo", "bar").Infof("hello world")
	clog.FromContext(ctx).With("bar", "baz").Infof("me again")
	clog.FromContext(ctx).With("baz", true).Infof("okay last one")

	clog.FromContext(ctx).Debug("hello debug")
	clog.FromContext(ctx).Info("hello info")
	clog.FromContext(ctx).Warn("hello warn")
	clog.FromContext(ctx).Error("hello error")

	fn(ctx)
}
