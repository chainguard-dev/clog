package slogtest_test

import (
	"testing"

	"github.com/chainguard-dev/clog"
	"github.com/chainguard-dev/clog/slogtest"
)

func TestSlogTest(t *testing.T) {
	ctx := slogtest.Context(t)

	clog.FromContext(ctx).With("foo", "bar").Infof("hello world")
	clog.FromContext(ctx).With("bar", "baz").Infof("me again")
	clog.FromContext(ctx).With("baz", true).Infof("okay last one")
}
