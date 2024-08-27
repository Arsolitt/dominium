package settings

import (
	"context"
	"fmt"
	"log/slog"
	"reflect"
	"sync"

	"github.com/arsolitt/dominium/lib/logger"
	"github.com/ilyakaznacheev/cleanenv"
	infisical "github.com/infisical/go-sdk"
)

type Settings struct {
	LogLevel              int    `env:"LOG_LEVEL"`
	InfisicalURL          string `env:"INFISICAL_URL"`
	InfisicalClientID     string `env:"INFISICAL_CLIENT_ID"`
	InfisicalClientSecret string `env:"INFISICAL_CLIENT_SECRET"`
	InfisicalProjectID    string `env:"INFISICAL_PROJECT_ID"`
	AppKey                string `infisical:"APP_KEY"`
}

var settings Settings
var once sync.Once

func Get() Settings {
	once.Do(func() {
		New()
	})
	return settings
}

func New() {
	tmp := Settings{}
	err := cleanenv.ReadEnv(&tmp)
	if err != nil {
		slog.Error("Failed to read env vars", "Error", err.Error())
		return
	}
	slog.Warn("Read settings from env vars")
	// err = readInfisical(&tmp)
	// if err != nil {
	// 	slog.Error("Failed to read remote config", "Error", err.Error(), "URL", tmp.InfisicalURL)
	// 	return
	// }
	// slog.Warn("Read settings from infisical")
	tags, err := parseStructTags(tmp)
	if err != nil {
		slog.Error("Failed to parse struct", "Error", err.Error())
	}
	slog.Warn("Parsed settings", "Tags", tags)
	for k, v := range tags {
		slog.Warn("Settings", "Key", tmp, "Value", v)
	}
	// fmt.Println(tags)
	settings = tmp
}

// func readInfisical(stg *Settings) error {
// 	client := infisical.NewInfisicalClient(infisical.Config{
// 		SiteUrl: stg.InfisicalURL,
// 	})

// 	_, err := client.Auth().UniversalAuthLogin(stg.InfisicalClientID, stg.InfisicalClientSecret)

// 	if err != nil {
// 		return logger.WrapError(context.TODO(), err)
// 	}

// 	appKey, err := client.Secrets().Retrieve(infisical.RetrieveSecretOptions{
// 		SecretKey:   "APP_KEY",
// 		Environment: "local",
// 		ProjectID:   stg.InfisicalProjectID,
// 		SecretPath:  "/",
// 	})
// 	if err != nil {
// 		return logger.WrapError(context.TODO(), err)
// 	}

// 	stg.AppKey = appKey.SecretValue
// 	return nil
// }

func readInfisical(stg *Settings) error {
	client := infisical.NewInfisicalClient(infisical.Config{
		SiteUrl: stg.InfisicalURL,
	})

	_, err := client.Auth().UniversalAuthLogin(stg.InfisicalClientID, stg.InfisicalClientSecret)

	if err != nil {
		return logger.WrapError(context.TODO(), err)
	}

	appKey, err := client.Secrets().Retrieve(infisical.RetrieveSecretOptions{
		SecretKey:   "APP_KEY",
		Environment: "local",
		ProjectID:   stg.InfisicalProjectID,
		SecretPath:  "/",
	})
	if err != nil {
		return logger.WrapError(context.TODO(), err)
	}

	stg.AppKey = appKey.SecretValue
	return nil
}

func parseStructTags(s any) (map[string]string, error) {
	typ := reflect.TypeOf(s)
	if typ.Kind() != reflect.Struct {
		return nil, fmt.Errorf("%s is not a struct", typ)
	}
	m := make(map[string]string)
	for i := 0; i < typ.NumField(); i++ {
		fld := typ.Field(i)
		if dbName := fld.Tag.Get("infisical"); dbName != "" {
			m[fld.Name] = dbName
		}
	}
	return m, nil
}
