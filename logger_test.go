package slogctx

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

	t.Run("slogctx.Info", func(t *testing.T) {
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
	want := fmt.Sprintf("github.com/wlynch/slogctx.%s", t.Name())
	if got.Source.Function != want {
		t.Errorf("want %v, got %v", want, got)
	}
}
