package slogctx

import (
	"bytes"
	"context"
	"encoding/json"
	"log/slog"
	"reflect"
	"testing"
)

func TestContextHandler(t *testing.T) {
	ctx := context.Background()
	ctx = WithValues(ctx, "foo", "bar")
	ctx2 := WithValues(ctx,
		"a", "b",
		"c", "d",
	)
	ctx = WithValues(ctx, "b", 1)

	for _, tc := range []struct {
		ctx  context.Context
		want map[string]any
	}{
		{ctx, map[string]any{
			"b":   float64(1),
			"foo": "bar",
		}},
		{ctx2, map[string]any{
			"foo": "bar",
			"a":   "b",
			"c":   "d",
		}},
	} {
		t.Run("", func(t *testing.T) {
			b := new(bytes.Buffer)
			log := slog.New(NewHandler(slog.NewJSONHandler(b, testopts)))
			log.InfoContext(tc.ctx, "")

			tc.want["level"] = "INFO"
			tc.want["msg"] = ""

			var got map[string]any
			if err := json.Unmarshal(b.Bytes(), &got); err != nil {
				t.Fatal(err)
			}

			if !reflect.DeepEqual(tc.want, got) {
				t.Errorf("want %v, got %v", tc.want, got)
			}
		})
	}
}
