package settings

import (
	"context"
	"fmt"
	"reflect"
	"sync"

	"github.com/arsolitt/dominium/lib/logger"
	infisical "github.com/infisical/go-sdk"
)

func infisicalSecret(field string, cfg InfisicalCreds, path string) (string, error) {
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

func readInfisicalConfig(stg interface{}, cfg InfisicalCreds, path string) error {
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
		wg.Add(1)
		go func(fld reflect.StructField, path string) {
			defer wg.Done()

			currentPath := path
			if newPath := fld.Tag.Get("infisical-path"); newPath != "" {
				currentPath = fmt.Sprintf("%s%s/", path, newPath)
			}

			if fld.Type.Kind() == reflect.Struct {
				err := readInfisicalConfig(val.Field(i).Addr().Interface(), cfg, currentPath)
				if err != nil {
					errCh <- err
				}
			} else if secretName := fld.Tag.Get("infisical"); secretName != "" {
				value, err := infisicalSecret(secretName, cfg, currentPath)
				if err != nil {
					errCh <- err
					return
				}

				fieldVal := val.Field(i)
				if !fieldVal.IsValid() || !fieldVal.CanSet() {
					errCh <- fmt.Errorf("field %s is not valid", fld.Name)
					return
				}

				mu.Lock()
				fieldVal.SetString(value)
				mu.Unlock()
			}
		}(fld, path)
	}

	wg.Wait()
	close(errCh)

	for err := range errCh {
		return logger.WrapError(context.TODO(), err)
	}

	return nil
}
