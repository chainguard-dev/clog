package slogtest_test

import (
	"context"

	"github.com/chainguard-dev/clog"
)

func fn(ctx context.Context) { clog.FromContext(ctx).With("foo", "bar").Infof("hello from fn") }
