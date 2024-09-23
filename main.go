package main

import (
	"context"
	"log/slog"
	"os"

	"github.com/arsolitt/dominium/lib/logger"
	"github.com/arsolitt/dominium/lib/settings"
	infisical "github.com/infisical/go-sdk"
)

func main() {
	// ctx := context.Background()
	logger.SetDefault(slog.LevelWarn, false)
	// stg := settings.Get()
	cfg := settings.GetCreds()
	err := test(cfg)

	if err != nil {
		slog.Error("Failed to read infisical vars", "Error", err.Error())
		os.Exit(1)
	}

	// for _, s := range secrets {
	// 	slog.Debug("Infisical secret", "key", s.SecretKey, "value", s.SecretValue)
	// }

	// ctx = logger.WithLogLevel(ctx, slog.Level(stg.LogLevel))
	// ctx = logger.WithLogValue(ctx, "app_config", stg)
	// slog.InfoContext(ctx, "App Started")
	// slog.DebugContext(ctx, "Debug Enabled")
}

func test(cfg settings.InfisicalCreds) error {
	client := infisical.NewInfisicalClient(infisical.Config{
		SiteUrl: cfg.InfisicalURL,
	})

	_, err := client.Auth().UniversalAuthLogin(cfg.InfisicalClientID, cfg.InfisicalClientSecret)

	if err != nil {
		return logger.WrapError(context.TODO(), err)
	}

	secrets, err := client.Secrets().List(infisical.ListSecretsOptions{
		ProjectID:          cfg.InfisicalProjectID,
		Environment:        cfg.Environment,
		SecretPath:         "/",
		AttachToProcessEnv: false,
	})
	if err != nil {
		return logger.WrapError(context.TODO(), err)
	}

	folders, err := client.Folders().List(infisical.ListFoldersOptions{
		ProjectID:   cfg.InfisicalProjectID,
		Environment: cfg.Environment,
		Path:        "/",
	})
	if err != nil {
		return logger.WrapError(context.TODO(), err)
	}

	// for _, s := range secrets {
	// 	slog.Warn("Infisical secret", "key", s.SecretKey, "value", s.SecretValue)
	// }
	slog.Warn("Infisical secrets", "secrets", secrets)

	slog.Warn("Infisical folders", "folders", folders)
	return nil
}
