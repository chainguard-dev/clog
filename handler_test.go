package clog

import (
	"bytes"
	"context"
	"encoding/json"
	"log"
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

// A Handler around a Handler, as FromContext builds when the default handler
// is already a Handler, writes each context value once.
func TestContextHandlerWrappedTwice(t *testing.T) {
	const once = `{"level":"INFO","msg":"","a":"b"}` + "\n"

	for _, tc := range []struct {
		name string
		log  func(ctx context.Context, h slog.Handler)
		want string
	}{{
		name: "NewHandler of a Handler",
		log: func(ctx context.Context, h slog.Handler) {
			slog.New(NewHandler(NewHandler(h))).InfoContext(ctx, "")
		},
		want: once,
	}, {
		name: "clog with a Handler as the default",
		log: func(ctx context.Context, h slog.Handler) {
			slog.SetDefault(slog.New(NewHandler(h)))
			InfoContext(ctx, "")
		},
		want: once,
	}, {
		name: "clog with a *Handler as the default",
		log: func(ctx context.Context, h slog.Handler) {
			hh := NewHandler(h)
			slog.SetDefault(slog.New(&hh))
			InfoContext(ctx, "")
		},
		want: once,
	}, {
		name: "slog with a Handler as the default",
		log: func(ctx context.Context, h slog.Handler) {
			slog.SetDefault(slog.New(NewHandler(h)))
			slog.InfoContext(ctx, "")
		},
		want: once,
	}, {
		name: "New(nil) with a Handler as the default",
		log: func(ctx context.Context, h slog.Handler) {
			slog.SetDefault(slog.New(NewHandler(h)))
			New(nil).InfoContext(ctx, "")
		},
		want: once,
	}, {
		name: "New(nil).With with a Handler as the default",
		log: func(ctx context.Context, h slog.Handler) {
			slog.SetDefault(slog.New(NewHandler(h)))
			New(nil).With("x", "y").InfoContext(ctx, "")
		},
		want: `{"level":"INFO","msg":"","x":"y","a":"b"}` + "\n",
	}} {
		t.Run(tc.name, func(t *testing.T) {
			previous, output, flags := slog.Default(), log.Writer(), log.Flags()
			t.Cleanup(func() {
				slog.SetDefault(previous)

				// slog.SetDefault also sends the log package's output to the
				// new handler, and restoring the previous default does not
				// undo that.
				log.SetOutput(output)
				log.SetFlags(flags)
			})

			b := new(bytes.Buffer)
			tc.log(WithValues(t.Context(), "a", "b"), slog.NewJSONHandler(b, testopts))

			if got := b.String(); got != tc.want {
				t.Errorf("got %q, want %q", got, tc.want)
			}
		})
	}
}
