package clog

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"runtime"
	"time"
)

// Logger implements a wrapper around [slog.Logger] that adds formatter functions (e.g. Infof, Errorf)
type Logger struct {
	ctx context.Context
	slog.Logger
}

// DefaultLogger returns a new logger that uses the default [slog.Logger].
func DefaultLogger() *Logger {
	return NewLogger(slog.Default())
}

// NewLogger returns a new logger that wraps the given [slog.Logger] with the default context.
func NewLogger(l *slog.Logger) *Logger {
	return NewLoggerWithContext(context.Background(), l)
}

// NewLoggerWithContext returns a new logger that wraps the given [slog.Logger].
func NewLoggerWithContext(ctx context.Context, l *slog.Logger) *Logger {
	if l == nil {
		l = slog.Default()
	}
	return &Logger{
		ctx:    ctx,
		Logger: *l,
	}
}

// New returns a new logger that wraps the given [slog.Handler].
func New(h slog.Handler) *Logger {
	return NewLogger(slog.New(h))
}

// New returns a new logger that wraps the given [slog.Handler].
func NewWithContext(ctx context.Context, h slog.Handler) *Logger {
	return NewLoggerWithContext(ctx, slog.New(h))
}

// With calls [Logger.With] on the default logger.
func With(args ...any) *Logger {
	return DefaultLogger().With(args...)
}

// With calls [Logger.With] on the logger.
func (l *Logger) With(args ...any) *Logger {
	return NewLoggerWithContext(l.context(), l.Logger.With(args...))
}

// WithGroup calls [Logger.WithGroup] on the default logger.
func (l *Logger) WithGroup(name string) *Logger {
	return NewLoggerWithContext(l.context(), l.Logger.WithGroup(name))
}

func (l *Logger) context() context.Context {
	if l.ctx == nil {
		return context.Background()
	}
	return l.ctx
}

// Infof logs at LevelInfo with the given format and arguments.
func (l *Logger) Infof(format string, args ...any) {
	wrapf(l.context(), l, slog.LevelInfo, format, args...)
}

// InfoContextf logs at LevelInfo with the given context, format and arguments.
func (l *Logger) InfoContextf(ctx context.Context, format string, args ...any) {
	wrapf(ctx, l, slog.LevelInfo, format, args...)
}

// Warnf logs at LevelWarn with the given format and arguments.
func (l *Logger) Warnf(format string, args ...any) {
	wrapf(l.context(), l, slog.LevelWarn, format, args...)
}

// WarnContextf logs at LevelWarn with the given context, format and arguments.
func (l *Logger) WarnContextf(ctx context.Context, format string, args ...any) {
	wrapf(ctx, l, slog.LevelWarn, format, args...)
}

// Errorf logs at LevelError with the given format and arguments.
func (l *Logger) Errorf(format string, args ...any) {
	wrapf(l.context(), l, slog.LevelError, format, args...)
}

// ErrorContextf logs at LevelError with the given context, format and arguments.
func (l *Logger) ErrorContextf(ctx context.Context, format string, args ...any) {
	wrapf(ctx, l, slog.LevelError, format, args...)
}

// Debugf logs at LevelDebug with the given format and arguments.
func (l *Logger) Debugf(format string, args ...any) {
	wrapf(l.context(), l, slog.LevelDebug, format, args...)
}

// DebugContextf logs at LevelDebug with the given context, format and arguments.
func (l *Logger) DebugContextf(ctx context.Context, format string, args ...any) {
	wrapf(ctx, l, slog.LevelDebug, format, args...)
}

// Fatalf logs at LevelError with the given format and arguments, then exits.
func (l *Logger) Fatalf(format string, args ...any) {
	wrapf(l.context(), l, slog.LevelError, format, args...)
	os.Exit(1)
}

// Fatal logs at LevelError with the given message, then exits.
func (l *Logger) Fatal(msg string, args ...any) {
	wrap(l.context(), l, slog.LevelError, msg, args...)
	os.Exit(1)
}

// FatalfContextf logs at LevelError with the given context, format and arguments, then exits.
func (l *Logger) FatalContextf(ctx context.Context, format string, args ...any) {
	wrapf(ctx, l, slog.LevelError, format, args...)
	os.Exit(1)
}

// FatalfContext logs at LevelError with the given context and message, then exits.
func (l *Logger) FatalContext(ctx context.Context, msg string, args ...any) {
	wrap(ctx, l, slog.LevelError, msg, args...)
	os.Exit(1)
}

// Base returns the underlying [slog.Logger].
func (l *Logger) Base() *slog.Logger {
	return &l.Logger
}

// Handler returns the underlying [slog.Handler].
func (l *Logger) Handler() slog.Handler {
	return l.Logger.Handler()
}

func wrap(ctx context.Context, logger *Logger, level slog.Level, msg string, args ...any) {
	if !logger.Handler().Enabled(ctx, level) {
		return
	}

	var pcs [1]uintptr
	runtime.Callers(3, pcs[:]) // skip [Callers, Infof, wrapf]
	r := slog.NewRecord(time.Now(), level, msg, pcs[0])
	r.Add(args...)
	_ = logger.Handler().Handle(ctx, r)
}

// wrapf is like wrap, but uses fmt.Sprintf to format the message.
// NOTE: args are passed to fmt.Sprintf, not as [slog.Attr].
func wrapf(ctx context.Context, logger *Logger, level slog.Level, format string, args ...any) {
	if !logger.Handler().Enabled(ctx, level) {
		return
	}

	var pcs [1]uintptr
	runtime.Callers(3, pcs[:]) // skip [Callers, Infof, wrapf]
	r := slog.NewRecord(time.Now(), level, fmt.Sprintf(format, args...), pcs[0])
	_ = logger.Handler().Handle(ctx, r)
}

type loggerKey struct{}

func WithLogger(ctx context.Context, logger *Logger) context.Context {
	return context.WithValue(ctx, loggerKey{}, logger.Logger)
}

func FromContext(ctx context.Context) *Logger {
	if logger, ok := ctx.Value(loggerKey{}).(slog.Logger); ok {
		return &Logger{
			ctx:    ctx,
			Logger: logger,
		}
	}
	return NewLoggerWithContext(ctx, slog.Default())
}
