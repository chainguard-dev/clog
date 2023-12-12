package slogctx

import (
	"context"
	"fmt"
	"log/slog"
)

func Info(ctx context.Context, args ...any) {
	slog.InfoContext(ctx, fmt.Sprint(args...))
}

func Infof(ctx context.Context, msg string, args ...any) {
	slog.InfoContext(ctx, fmt.Sprintf(msg, args...))
}

func Warn(ctx context.Context, args ...any) {
	slog.WarnContext(ctx, fmt.Sprint(args...))
}

func Warnf(ctx context.Context, msg string, args ...any) {
	slog.WarnContext(ctx, fmt.Sprintf(msg, args...))
}

func Error(ctx context.Context, args ...any) {
	slog.ErrorContext(ctx, fmt.Sprint(args...))
}

func Errorf(ctx context.Context, msg string, args ...any) {
	slog.ErrorContext(ctx, fmt.Sprintf(msg, args...))
}

func Debug(ctx context.Context, args ...any) {
	slog.DebugContext(ctx, fmt.Sprint(args...))
}

func Debugf(ctx context.Context, msg string, args ...any) {
	slog.DebugContext(ctx, fmt.Sprintf(msg, args...))
}
