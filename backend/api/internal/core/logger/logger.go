package logger

import (
	"context"
	"log/slog"
	"os"
)

type Logger interface {
	Info(ctx context.Context, message string, attrs ...Attr)
	Warn(ctx context.Context, message string, attrs ...Attr)
	Error(ctx context.Context, message string, attrs ...Attr)
}

type Attr = slog.Attr

type SlogLogger struct {
	logger *slog.Logger
}

func New(env string) Logger {
	level := slog.LevelInfo
	if env == "development" {
		level = slog.LevelDebug
	}

	handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: level})
	return SlogLogger{logger: slog.New(handler)}
}

func (l SlogLogger) Info(ctx context.Context, message string, attrs ...Attr) {
	l.logger.LogAttrs(ctx, slog.LevelInfo, message, attrs...)
}

func (l SlogLogger) Warn(ctx context.Context, message string, attrs ...Attr) {
	l.logger.LogAttrs(ctx, slog.LevelWarn, message, attrs...)
}

func (l SlogLogger) Error(ctx context.Context, message string, attrs ...Attr) {
	l.logger.LogAttrs(ctx, slog.LevelError, message, attrs...)
}
