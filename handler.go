package slogctx

import (
	"context"
	"log/slog"
)

var (
	ctxKey = struct{}{}
)

type ctxVal map[string]any

// With returns a new context with the given values.
// Values are expected to be key-value pairs, where the key is a string.
// e.g. With(ctx, "foo", "bar", "baz", 1)
// If a value already exists, it is overwritten.
// If an odd number of arguments are provided, With panics.
func With(ctx context.Context, args ...any) context.Context {
	if len(args)%2 != 0 {
		panic("non-even number of arguments")
	}

	values := ctxVal{}

	// Copy existing values
	for k, v := range get(ctx) {
		values[k] = v
	}

	for i := 0; i < len(args); i++ {
		key, ok := args[i].(string)
		if !ok {
			panic("non-string key")
		}
		i++
		if i >= len(args) {
			break
		}
		value := args[i]
		values[key] = value
	}
	return context.WithValue(ctx, ctxKey, values)
}

func get(ctx context.Context) ctxVal {
	if value, ok := ctx.Value(ctxKey).(ctxVal); ok {
		return value
	}
	return nil
}

// Handler is a slog.Handler that adds context values to the log record.
// Values are added via [With].
type Handler struct {
	h slog.Handler
}

// NewHandler configures a new context aware slog handler.
// If h is nil, the default slog handler is used.
func NewHandler(h slog.Handler) Handler {
	return Handler{h}
}

func (h Handler) inner() slog.Handler {
	if h.h == nil {
		return slog.Default().Handler()
	}
	return h.h
}

func (h Handler) Enabled(ctx context.Context, level slog.Level) bool {
	return true
}

func (h Handler) Handle(ctx context.Context, r slog.Record) error {
	values := get(ctx)
	for k, v := range values {
		r.Add(k, v)
	}
	return h.inner().Handle(ctx, r)
}

func (h Handler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return Handler{h.inner().WithAttrs(attrs)}
}

func (h Handler) WithGroup(name string) slog.Handler {
	return Handler{h.inner().WithGroup(name)}
}
