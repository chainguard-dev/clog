package slogctx

import (
	"bytes"
	"context"
	"log/slog"
	"testing"
)

func TestContextHandler(t *testing.T) {
	ctx := context.Background()
	ctx = With(ctx, "foo", "bar")
	ctx2 := With(ctx,
		"a", "b",
		"c", "d",
	)
	ctx = With(ctx, "b", 1)

	for _, tc := range []struct {
		ctx  context.Context
		want string
	}{
		{ctx, "foo=bar b=1"},
		{ctx2, "foo=bar a=b c=d"},
	} {
		t.Run("", func(t *testing.T) {
			b := new(bytes.Buffer)
			log := slog.New(NewHandler(slog.NewTextHandler(b, nil)))
			log.InfoContext(tc.ctx, "hello world")
		})
	}
}
