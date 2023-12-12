package main

import (
	"testing"

	"github.com/wlynch/slogctx"
	"github.com/wlynch/slogctx/slogtest"
)

func TestFoo(t *testing.T) {
	ctx := slogtest.TestContextWithLogger(t)

	for _, tc := range []string{"a", "b"} {
		t.Run(tc, func(t *testing.T) {
			slogctx.FromContext(ctx).Infof("hello world")
		})
	}
}
