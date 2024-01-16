package clog

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"reflect"
	"testing"
)

var (
	testopts = &slog.HandlerOptions{
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			// Ignore time to make testing easier.
			if a.Key == "time" {
				return slog.Attr{}
			}
			return a
		},
	}
)

func TestLogger(t *testing.T) {
	ctx := context.Background()
	b := new(bytes.Buffer)
	base := slog.New(NewHandler(slog.NewJSONHandler(b, testopts)))
	log := NewLogger(base).With("a", "b")
	log.InfoContext(ctx, "")
	t.Log(b.String())

	want := map[string]any{
		"level": "INFO",
		"msg":   "",
		"a":     "b",
	}

	var got map[string]any
	if err := json.Unmarshal(b.Bytes(), &got); err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(want, got) {
		t.Errorf("want %v, got %v", want, got)
	}
}

func TestLoggerNilBase(t *testing.T) {
	log := NewLogger(nil)
	log.Info("")
}

func TestLoggerFromContext(t *testing.T) {
	b := new(bytes.Buffer)
	base := slog.New(NewHandler(slog.NewJSONHandler(b, testopts)))
	log := NewLogger(base).With("a", "b")

	ctx := WithLogger(context.Background(), log)
	FromContext(ctx).Info("")

	want := map[string]any{
		"level": "INFO",
		"msg":   "",
		"a":     "b",
	}

	t.Run("FromContext.Info", func(t *testing.T) {
		var got map[string]any
		if err := json.Unmarshal(b.Bytes(), &got); err != nil {
			t.Fatal(err)
		}
		if !reflect.DeepEqual(want, got) {
			t.Errorf("want %v, got %v", want, got)
		}
	})

	b.Reset()

	t.Run("clog.Info", func(t *testing.T) {
		InfoContext(ctx, "")
		var got map[string]any
		if err := json.Unmarshal(b.Bytes(), &got); err != nil {
			t.Fatal(err)
		}
		if !reflect.DeepEqual(want, got) {
			t.Errorf("want %v, got %v", want, got)
		}
	})
}

func TestLoggerPC(t *testing.T) {
	b := new(bytes.Buffer)
	log := NewLogger(slog.New(NewHandler(slog.NewJSONHandler(b, &slog.HandlerOptions{
		AddSource:   true,
		ReplaceAttr: testopts.ReplaceAttr,
	}))))

	log.Info("")
	t.Log(b.String())

	var got struct {
		Source struct {
			File     string `json:"file"`
			Function string `json:"function"`
		} `json:"source"`
	}
	if err := json.Unmarshal(b.Bytes(), &got); err != nil {
		t.Fatal(err)
	}
	// Knowing that the PC is from this test is good enough.
	want := fmt.Sprintf("github.com/chainguard-dev/clog.%s", t.Name())
	if got.Source.Function != want {
		t.Errorf("want %v, got %v", want, got)
	}
}

func TestContext(t *testing.T) {
	// Stuff some data in the context, to check that it gets passed through to the slog handler.
	h := &testHandler{t: t}
	ctx := context.Background()
	ctx = context.WithValue(ctx, logKey{}, "value")
	ctx = WithLogger(ctx, NewLogger(slog.New(h)))

	// Calling InfoContext uses the context passed explicitly.
	FromContext(ctx).InfoContext(ctx, "with explicit context")

	// Calling context-less Info and Infof uses the context from FromContext.
	FromContext(ctx).Info("with implicit context")
	FromContext(ctx).Infof("with implicit context %q", "and format")
}

type logKey struct{}

type testHandler struct{ t *testing.T }

func (_ *testHandler) Enabled(context.Context, slog.Level) bool { return true }
func (t *testHandler) WithAttrs(attrs []slog.Attr) slog.Handler { return t }
func (t *testHandler) WithGroup(name string) slog.Handler       { return t }
func (t *testHandler) Handle(ctx context.Context, r slog.Record) error {
	if got := ctx.Value(logKey{}); got != "value" {
		t.t.Errorf("%s: expected value, got %v", r.Message, got)
	}
	return nil
}
