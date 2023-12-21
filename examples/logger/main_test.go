package main

import (
	"testing"

	"github.com/chainguard-dev/slogctx"
	"github.com/chainguard-dev/slogctx/slogtest"
)

func TestFoo(t *testing.T) {
	ctx := slogtest.TestContextWithLogger(t)

	for _, tc := range []string{"a", "b"} {
		t.Run(tc, func(t *testing.T) {
			slogctx.FromContext(ctx).Infof("hello world")
		})
	}
}
