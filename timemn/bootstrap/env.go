package bootstrap

import (
	"github.com/go-viper/mapstructure/v2"
	"github.com/knadh/koanf/parsers/dotenv"
	"github.com/knadh/koanf/providers/env"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"
	"github.com/nhuongmh/cvfs/timemn/internal/logger"
)

type Env struct {
	AppMode        string `mapstructure:"APP_MODE"`
	ContextTimeout int    `mapstructure:"CONTEXT_TIMEOUT"`
	ServerAddress  string `mapstructure:"SERVER_ADDRESS"`
	DBHost         string `mapstructure:"DB_HOST"`
	DBPort         string `mapstructure:"DB_PORT"`
	DBUser         string `mapstructure:"DB_USERNAME"`
	DBPassword     string `mapstructure:"DB_PASSWORD"`
	DBName         string `mapstructure:"DB_NAME"`
	DBSchema       string `mapstructure:"DB_SCHEMA"`
}

func NewEnv() *Env {
	logger.Log.Info().Msg("Reading ENV")
	envData := Env{}
	var k = koanf.New(".")
	if err := k.Load(file.Provider(".env"), dotenv.Parser()); err != nil {
		logger.Log.Fatal().Err(err).Msg("Failed load .env config file")
	}
	k.Load(env.Provider("", ".", nil), nil)
	if err := mapstructure.WeakDecode(k.All(), &envData); err != nil {
		logger.Log.Fatal().Err(err).Msg("Failed parse Env variable")
	}

	if envData.AppMode == "dev" {
		logger.Log.Info().Msg("The application is runnning in development mode")
	}

	return &envData
}
