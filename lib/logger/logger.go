package logger

import (
	"context"
	"log/slog"
	"os"
	"sync"
)

var once sync.Once

func SetDefault(level slog.Level, source bool) {
	once.Do(func() {
		handler := slog.Handler(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level:     level,
			AddSource: source,
		}))
		handler = NewMiddleware(handler)
		slog.SetDefault(slog.New(handler))
	})
}

func WithLogValue(ctx context.Context, entryKey string, value any) context.Context {
	if c, ok := ctx.Value(dataKey).(logData); ok {
		c[entryKey] = value
		return context.WithValue(ctx, dataKey, c)
	}
	return context.WithValue(ctx, dataKey, logData{entryKey: value})
}

func WithLogLevel(ctx context.Context, value slog.Level) context.Context {
	return context.WithValue(ctx, levelKey, value)
}
