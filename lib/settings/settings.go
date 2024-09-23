package settings

import (
	"log/slog"
	"os"
	"sync"

	"github.com/ilyakaznacheev/cleanenv"
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
	Aboba       string   `infisical:"ABOBA" infisical-path:"aboba"`
	Buba        Buba     `infisical-path:"features/buba"`
	Features    Features `infisical-path:"features"`
}

type Features struct {
	UserLogLevel string `infisical:"USER_LOG_LEVEL"`
}

type Buba struct {
	BamBam string `infisical:"BAM_BAM"`
}

type InfisicalCreds struct {
	Environment           string `env:"ENVIRONMENT"`
	InfisicalURL          string `env:"INFISICAL_URL"`
	InfisicalClientID     string `env:"INFISICAL_CLIENT_ID"`
	InfisicalClientSecret string `env:"INFISICAL_CLIENT_SECRET"`
	InfisicalProjectID    string `env:"INFISICAL_PROJECT_ID"`
}

var settings Settings
var cfg InfisicalCreds
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
		err = readInfisicalConfig(&settings, cfg, "/")
		if err != nil {
			slog.Error("Failed to hydrate settings", "Error", err.Error())
			os.Exit(1)
		}
	})
	return settings
}

func GetCreds() InfisicalCreds {
	once.Do(func() {
		slog.Warn("Read infisical creds from env vars")
		err := cleanenv.ReadEnv(&cfg)
		if err != nil {
			slog.Error("Failed to read env vars", "Error", err.Error())
			os.Exit(1)
		}

	})
	return cfg
}
