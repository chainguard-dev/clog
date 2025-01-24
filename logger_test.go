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

func TestWith(t *testing.T) {
	ctx := context.WithValue(context.Background(), "test", "test")
	log := NewLoggerWithContext(ctx, nil)
	withed := log.With("a", "b")
	if want := withed.ctx; want != ctx {
		t.Errorf("want %v, got %v", want, ctx)
	}
	withed = log.WithGroup("a")
	if want := withed.ctx; want != ctx {
		t.Errorf("want %v, got %v", want, ctx)
	}
}

func TestDefaultHandler(t *testing.T) {
	old := slog.Default()
	defer func() {
		slog.SetDefault(old)
	}()

	b := new(bytes.Buffer)
	slog.SetDefault(slog.New(slog.NewJSONHandler(b, testopts)))

	t.Run("Info", func(t *testing.T) {
		FromContext(WithValues(context.Background(), "a", "b")).Info("")
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
	})

	b.Reset()

	t.Run("InfoContext", func(t *testing.T) {
		// Set logger with original value
		ctx := WithValues(context.Background(), "a", "b")
		logger := FromContext(ctx)

		// Override value in request context - we expect this to overwrite the original value set in the logger
		logger.InfoContext(WithValues(ctx, "a", "c"), "")

		want := map[string]any{
			"level": "INFO",
			"msg":   "",
			"a":     "c",
		}
		var got map[string]any
		if err := json.Unmarshal(b.Bytes(), &got); err != nil {
			t.Fatal(err)
		}
		if !reflect.DeepEqual(want, got) {
			t.Errorf("want %v, got %v", want, got)
		}
	})
}
