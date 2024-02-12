package slogtest

import (
	"testing"

	"github.com/chainguard-dev/clog"
)

func TestSlogTest(t *testing.T) {
	ctx := TestContextWithLogger(t)

	clog.FromContext(ctx).With("foo", "bar").Infof("hello world")
	clog.FromContext(ctx).With("bar", "baz").Infof("me again")
	clog.FromContext(ctx).With("baz", true).Infof("okay last one")
}
