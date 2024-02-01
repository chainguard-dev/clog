package gcp

import (
	"context"
	"log/slog"
	"os"
)

// LevelCritical is an extra log level supported by Cloud Logging.
const LevelCritical = slog.Level(12)

// Handler that outputs JSON understood by the structured log agent.
// See https://cloud.google.com/logging/docs/agent/logging/configuration#special-fields
type Handler struct{ handler slog.Handler }

func NewHandler() *Handler {
	return &Handler{handler: slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{
		AddSource: true,
		Level:     slog.LevelDebug,
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			if a.Key == slog.MessageKey {
				a.Key = "message"
			} else if a.Key == slog.SourceKey {
				a.Key = "logging.googleapis.com/sourceLocation"
			} else if a.Key == slog.LevelKey {
				a.Key = "severity"
				level := a.Value.Any().(slog.Level)
				if level == LevelCritical {
					a.Value = slog.StringValue("CRITICAL")
				}
			}
			return a
		},
	})}
}

func (h *Handler) Enabled(ctx context.Context, level slog.Level) bool {
	return h.handler.Enabled(ctx, level)
}

func (h *Handler) Handle(ctx context.Context, rec slog.Record) error {
	if trace := traceFromContext(ctx); trace != "" {
		rec = rec.Clone()
		// Add trace ID	to the record so it is correlated with the request log
		// See https://cloud.google.com/trace/docs/trace-log-integration
		rec.Add("logging.googleapis.com/trace", slog.StringValue(trace))
	}

	return h.handler.Handle(ctx, rec)
}

func (h *Handler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &Handler{handler: h.handler.WithAttrs(attrs)}
}

func (h *Handler) WithGroup(name string) slog.Handler {
	return &Handler{handler: h.handler.WithGroup(name)}
}
