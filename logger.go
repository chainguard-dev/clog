package slogctx

import (
	"context"
	"fmt"
	"log/slog"
	"runtime"
	"time"
)

// Logger implements a wrapper around [slog.Logger] that adds formatter functions (e.g. Infof, Errorf)
type Logger struct {
	slog.Logger
}

// DefaultLogger returns a new logger that uses the default [slog.Logger].
func DefaultLogger() *Logger {
	return NewLogger(slog.Default())
}

// NewLogger returns a new logger that wraps the given [slog.Logger].
func NewLogger(l *slog.Logger) *Logger {
	if l == nil {
		l = slog.Default()
	}
	return &Logger{Logger: *l}
}

// New returns a new logger that wraps the given [slog.Handler].
func New(h slog.Handler) *Logger {
	return NewLogger(slog.New(h))
}

// With calls [Logger.With] on the default logger.
func With(args ...any) *Logger {
	return DefaultLogger().With(args...)
}

// With calls [Logger.With] on the logger.
func (l *Logger) With(args ...any) *Logger {
	return NewLogger(l.Logger.With(args...))
}

// WithGroup calls [Logger.WithGroup] on the default logger.
func (l *Logger) WithGroup(name string) *Logger {
	return NewLogger(l.Logger.WithGroup(name))
}

// Infof logs at LevelInfo with the given format and arguments.
func (l *Logger) Infof(format string, args ...any) {
	wrapf(context.Background(), l, slog.LevelInfo, format, args...)
}

// InfoContextf logs at LevelInfo with the given context, format and arguments.
func (l *Logger) InfoContextf(ctx context.Context, format string, args ...any) {
	wrapf(ctx, l, slog.LevelInfo, format, args...)
}

// Warnf logs at LevelWarn with the given format and arguments.
func (l *Logger) Warnf(format string, args ...any) {
	wrapf(context.Background(), l, slog.LevelWarn, format, args...)
}

// WarnContextf logs at LevelWarn with the given context, format and arguments.
func (l *Logger) WarnContextf(ctx context.Context, format string, args ...any) {
	wrapf(ctx, l, slog.LevelWarn, format, args...)
}

// Errorf logs at LevelError with the given format and arguments.
func (l *Logger) Errorf(format string, args ...any) {
	wrapf(context.Background(), l, slog.LevelError, format, args...)
}

// ErrorContextf logs at LevelError with the given context, format and arguments.
func (l *Logger) ErrorContextf(ctx context.Context, format string, args ...any) {
	wrapf(ctx, l, slog.LevelError, format, args...)
}

// Debugf logs at LevelDebug with the given format and arguments.
func (l *Logger) Debugf(format string, args ...any) {
	wrapf(context.Background(), l, slog.LevelDebug, format, args...)
}

// DebugContextf logs at LevelDebug with the given context, format and arguments.
func (l *Logger) DebugContextf(ctx context.Context, format string, args ...any) {
	wrapf(ctx, l, slog.LevelDebug, format, args...)
}

// Base returns the underlying [slog.Logger].
func (l *Logger) Base() *slog.Logger {
	return &l.Logger
}

func (l *Logger) Handler() slog.Handler {
	return l.Logger.Handler()
}

func wrap(ctx context.Context, logger *Logger, level slog.Level, msg string, args ...any) {
	var pcs [1]uintptr
	runtime.Callers(3, pcs[:]) // skip [Callers, Infof, wrapf]
	r := slog.NewRecord(time.Now(), level, msg, pcs[0])
	r.Add(args...)
	_ = logger.Handler().Handle(ctx, r)
}

// wrapf is like wrap, but uses fmt.Sprintf to format the message.
// NOTE: args are passed to fmt.Sprintf, not as [slog.Attr].
func wrapf(ctx context.Context, logger *Logger, level slog.Level, format string, args ...any) {
	var pcs [1]uintptr
	runtime.Callers(3, pcs[:]) // skip [Callers, Infof, wrapf]
	r := slog.NewRecord(time.Now(), level, fmt.Sprintf(format, args...), pcs[0])
	_ = logger.Handler().Handle(ctx, r)
}

type loggerKey struct{}

func WithLogger(ctx context.Context, logger *Logger) context.Context {
	return context.WithValue(ctx, loggerKey{}, logger)
}

func FromContext(ctx context.Context) *Logger {
	if logger, ok := ctx.Value(loggerKey{}).(*Logger); ok {
		return logger
	}
	return DefaultLogger()
}
