package settings

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"reflect"
	"sync"

	"github.com/arsolitt/dominium/lib/logger"
	"github.com/ilyakaznacheev/cleanenv"
	infisical "github.com/infisical/go-sdk"
)

type Settings struct {
	LogLevel    int      `env:"LOG_LEVEL"`
	Environment string   `env:"ENVIRONMENT"`
	AppKey      string   `infisical:"APP_KEY"`
	DbHost      string   `infisical:"DB_HOST"`
	DbPort      string   `infisical:"DB_PORT"`
	DbUser      string   `infisical:"DB_USER"`
	DbPassword  string   `infisical:"DB_PASSWORD"`
	DbName      string   `infisical:"DB_NAME"`
	CacheHost   string   `infisical:"CACHE_HOST"`
	CachePort   string   `infisical:"CACHE_PORT"`
	Aboba       string   `infisical:"ABOBA" infisical-path:"/aboba"`
	Features    Features `infisical-path:"/features"`
}

type Features struct {
	UserLogLevel string `infisical:"USER_LOG_LEVEL"`
}

type infisicalCreds struct {
	Environment           string `env:"ENVIRONMENT"`
	InfisicalURL          string `env:"INFISICAL_URL"`
	InfisicalClientID     string `env:"INFISICAL_CLIENT_ID"`
	InfisicalClientSecret string `env:"INFISICAL_CLIENT_SECRET"`
	InfisicalProjectID    string `env:"INFISICAL_PROJECT_ID"`
}

var settings Settings
var cfg infisicalCreds
var once sync.Once

func Get() Settings {
	once.Do(func() {

		slog.Warn("Read settings from env vars")
		err := cleanenv.ReadEnv(&settings)
		if err != nil {
			slog.Error("Failed to read env vars", "Error", err.Error())
			os.Exit(1)
		}

		slog.Warn("Read infisical creds from env vars")
		err = cleanenv.ReadEnv(&cfg)
		if err != nil {
			slog.Error("Failed to read env vars", "Error", err.Error())
			os.Exit(1)
		}

		slog.Warn("Read settings from infisical")
		err = hydrateInfisicalSettings(&settings, cfg, "/")
		if err != nil {
			slog.Error("Failed to hydrate settings", "Error", err.Error())
			os.Exit(1)
		}
	})
	return settings
}

func readInfisical(field string, cfg infisicalCreds, path string) (string, error) {
	client := infisical.NewInfisicalClient(infisical.Config{
		SiteUrl: cfg.InfisicalURL,
	})

	_, err := client.Auth().UniversalAuthLogin(cfg.InfisicalClientID, cfg.InfisicalClientSecret)

	if err != nil {
		return "", logger.WrapError(context.TODO(), err)
	}

	secret, err := client.Secrets().Retrieve(infisical.RetrieveSecretOptions{
		SecretKey:   field,
		Environment: cfg.Environment,
		ProjectID:   cfg.InfisicalProjectID,
		SecretPath:  path,
	})
	if err != nil {
		return "", logger.WrapError(context.TODO(), err)
	}

	if secret.SecretValue == "" {
		return "", logger.WrapError(context.TODO(), fmt.Errorf("secret %s is empty", field))
	}
	return secret.SecretValue, nil
}

// func hydrateInfisicalSettings(stg *Settings) error {
// 	typ := reflect.TypeOf(*stg)
// 	if typ.Kind() != reflect.Struct {
// 		return logger.WrapError(context.TODO(), fmt.Errorf("%s is not a struct", typ))
// 	}

// 	val := reflect.ValueOf(stg).Elem()
// 	for i := 0; i < typ.NumField(); i++ {
// 		fld := typ.Field(i)
// 		secretName := fld.Tag.Get("infisical")
// 		if secretName == "" {
// 			continue
// 		}

// 		value, err := readInfisical(secretName, stg)
// 		if err != nil {
// 			return logger.WrapError(context.TODO(), err)
// 		}

// 		fieldVal := val.Field(i)
// 		if !fieldVal.IsValid() || !fieldVal.CanSet() {
// 			continue
// 		}

// 		fieldVal.SetString(value)
// 	}

// 	return nil
// }

func hydrateInfisicalSettings(stg interface{}, cfg infisicalCreds, path string) error {
	typ := reflect.TypeOf(stg).Elem()
	if typ.Kind() != reflect.Struct {
		return logger.WrapError(context.TODO(), fmt.Errorf("%s is not a struct", typ))
	}

	val := reflect.ValueOf(stg).Elem()
	var wg sync.WaitGroup
	var mu sync.Mutex
	errCh := make(chan error, typ.NumField())

	for i := 0; i < typ.NumField(); i++ {
		fld := typ.Field(i)
		secretName := fld.Tag.Get("infisical")
		if secretName != "" {

			wg.Add(1)
			go func(i int, secretName string) {
				defer wg.Done()
				newPath := fld.Tag.Get("infisical-path")
				tmp := path
				if newPath != "" {
					tmp = newPath
				}
				value, err := readInfisical(secretName, cfg, tmp)
				if err != nil {
					errCh <- err
					return
				}

				fieldVal := val.Field(i)
				if !fieldVal.IsValid() || !fieldVal.CanSet() {
					return
				}

				mu.Lock()
				fieldVal.SetString(value)
				mu.Unlock()
			}(i, secretName)
		} else if fld.Type.Kind() == reflect.Struct {
			wg.Add(1)
			go func(i int) {
				defer wg.Done()
				newPath := fld.Tag.Get("infisical-path")
				if newPath != "" {
					err := hydrateInfisicalSettings(val.Field(i).Addr().Interface(), cfg, newPath)
					if err != nil {
						errCh <- err
					}
				} else {
					err := hydrateInfisicalSettings(val.Field(i).Addr().Interface(), cfg, path)
					if err != nil {
						errCh <- err
					}
				}
			}(i)
		}
	}

	wg.Wait()
	close(errCh)

	for err := range errCh {
		return logger.WrapError(context.TODO(), err)
	}

	return nil
}
