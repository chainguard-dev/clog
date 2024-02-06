package clog

import (
	"context"
	"log/slog"
	"os"
)

// Info calls Info on the default logger.
func Info(msg string, args ...any) {
	wrap(context.Background(), DefaultLogger(), slog.LevelInfo, msg, args...)
}

// InfoContext calls InfoContext on the context logger.
// If a Logger is found in the context, it will be used.
func InfoContext(ctx context.Context, msg string, args ...any) {
	wrap(ctx, FromContext(ctx), slog.LevelInfo, msg, args...)
}

// Infof calls Infof on the default logger.
func Infof(format string, args ...any) {
	wrapf(context.Background(), DefaultLogger(), slog.LevelInfo, format, args...)
}

// InfoContextf calls InfoContextf on the context logger.
// If a Logger is found in the context, it will be used.
func InfoContextf(ctx context.Context, format string, args ...any) {
	wrapf(ctx, FromContext(ctx), slog.LevelInfo, format, args...)
}

// Warn calls Warn on the default logger.
func Warn(msg string, args ...any) {
	wrap(context.Background(), DefaultLogger(), slog.LevelWarn, msg, args...)
}

// WarnContext calls WarnContext on the context logger.
// If a Logger is found in the context, it will be used.
func WarnContext(ctx context.Context, msg string, args ...any) {
	wrap(ctx, FromContext(ctx), slog.LevelWarn, msg, args...)
}

// Warnf calls Warnf on the default logger.
func Warnf(format string, args ...any) {
	wrapf(context.Background(), DefaultLogger(), slog.LevelWarn, format, args...)
}

// WarnContextf calls WarnContextf on the context logger.
// If a Logger is found in the context, it will be used.
func WarnContextf(ctx context.Context, format string, args ...any) {
	wrapf(ctx, FromContext(ctx), slog.LevelWarn, format, args...)
}

// Error calls Error on the default logger.
func Error(msg string, args ...any) {
	wrap(context.Background(), DefaultLogger(), slog.LevelError, msg, args...)
}

// ErrorContext calls ErrorContext on the context logger.
func ErrorContext(ctx context.Context, msg string, args ...any) {
	wrap(ctx, FromContext(ctx), slog.LevelError, msg, args...)
}

// Errorf calls Errorf on the default logger.
func Errorf(format string, args ...any) {
	wrapf(context.Background(), DefaultLogger(), slog.LevelError, format, args...)
}

// ErrorContextf calls ErrorContextf on the context logger.
func ErrorContextf(ctx context.Context, format string, args ...any) {
	wrapf(ctx, FromContext(ctx), slog.LevelError, format, args...)
}

// Debug calls Debug on the default logger.
func Debug(msg string, args ...any) {
	wrap(context.Background(), DefaultLogger(), slog.LevelDebug, msg, args...)
}

// DebugContext calls DebugContext on the context logger.
func DebugContext(ctx context.Context, msg string, args ...any) {
	wrap(ctx, FromContext(ctx), slog.LevelDebug, msg, args...)
}

// Debugf calls Debugf on the default logger.
func Debugf(format string, args ...any) {
	wrapf(context.Background(), DefaultLogger(), slog.LevelDebug, format, args...)
}

// DebugContextf calls DebugContextf on the context logger.
// If a Logger is found in the context, it will be used.
func DebugContextf(ctx context.Context, format string, args ...any) {
	wrapf(ctx, FromContext(ctx), slog.LevelDebug, format, args...)
}

// Fatal calls Error on the default logger, then exits.
func Fatal(msg string, args ...any) {
	wrap(context.Background(), DefaultLogger(), slog.LevelError, msg, args...)
	os.Exit(1)
}

// FatalContext calls ErrorContext on the context logger, then exits.
func FatalContext(ctx context.Context, msg string, args ...any) {
	wrap(ctx, FromContext(ctx), slog.LevelError, msg, args...)
	os.Exit(1)
}

// Fatalf calls Errorf on the default logger, then exits.
func Fatalf(format string, args ...any) {
	wrapf(context.Background(), DefaultLogger(), slog.LevelError, format, args...)
	os.Exit(1)
}

// FatalContextf calls ErrorContextf on the context logger, then exits.
func FatalContextf(ctx context.Context, format string, args ...any) {
	wrapf(ctx, FromContext(ctx), slog.LevelError, format, args...)
	os.Exit(1)
}
