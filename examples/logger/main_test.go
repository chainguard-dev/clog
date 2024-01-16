package main

import (
	"testing"

	"github.com/chainguard-dev/clog"
	"github.com/chainguard-dev/clog/slogtest"
)

func TestFoo(t *testing.T) {
	ctx := slogtest.TestContextWithLogger(t)

	for _, tc := range []string{"a", "b"} {
		t.Run(tc, func(t *testing.T) {
			clog.FromContext(ctx).Infof("hello world")
		})
	}
}
