package main

import (
	"context"
	"log/slog"

	"github.com/arsolitt/dominium/lib/logger"
	"github.com/arsolitt/dominium/lib/settings"
)

func main() {
	ctx := context.Background()
	logger.SetDefault(slog.LevelWarn, false)
	stg := settings.Get()

	ctx = logger.WithLogLevel(ctx, slog.Level(stg.LogLevel))
	ctx = logger.WithLogValue(ctx, "app_config", stg)
	slog.InfoContext(ctx, "App Started")
	slog.DebugContext(ctx, "Debug Enabled")
}
