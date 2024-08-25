package main

import (
	"context"
	"errors"
	"log/slog"

	"github.com/arsolitt/dominium/lib/logger"
)

func main() {
	logger.InitLogging()
	ctx := context.Background()

	reqID := "123121"
	ctx = logger.WithLogValue(ctx, logger.RequestIDField, reqID)
	slog.InfoContext(ctx, "New request")

	userId := "42"
	ctx = logger.WithLogValue(ctx, logger.UserIDField, userId)
	slog.InfoContext(ctx, "Processing user")

	instanceId := "228"
	ctx = logger.WithLogValue(ctx, logger.InstanceIDField, instanceId)
	slog.InfoContext(ctx, "Processing instance")

	err := errors.New("some error")
	logger.WrapError(ctx, err)
	slog.ErrorContext(ctx, "Processing request error")

	slog.DebugContext(ctx, "Debug message")

	slog.InfoContext(ctx, "Done")
}
