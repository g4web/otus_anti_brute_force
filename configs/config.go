package configs

import (
	"log"
	"time"

	"github.com/spf13/viper"
)

type configMap struct {
	IPTimeLimit                  string `mapstructure:"TIME_LIMIT_IP"`
	IPMaxCountForTimeLimit       int    `mapstructure:"MAX_COUNT_FOR_TIME_LIMIT_IP"`
	LoginTimeLimit               string `mapstructure:"TIME_LIMIT_LOGIN"`
	LoginMaxCountForTimeLimit    int    `mapstructure:"MAX_COUNT_FOR_TIME_LIMIT_LOGIN"`
	PasswordTimeLimit            string `mapstructure:"TIME_LIMIT_PASSWORD"`
	PasswordMaxCountForTimeLimit int    `mapstructure:"MAX_COUNT_FOR_TIME_LIMIT_PASSWORD"`
	GrpcHost                     string `mapstructure:"GRPC_HOST"`
	GrpcPort                     string `mapstructure:"GRPC_PORT"`
	CleanupPeriod                string `mapstructure:"CLEANUP_PERIOD"`
}

type Config struct {
	IPTimeLimit                  time.Duration
	IPMaxCountForTimeLimit       int
	LoginTimeLimit               time.Duration
	LoginMaxCountForTimeLimit    int
	PasswordTimeLimit            time.Duration
	PasswordMaxCountForTimeLimit int
	GrpcHost                     string
	GrpcPort                     string
	CleanUpPeriod                time.Duration
}

func NewConfig(filePath string) (*Config, error) {
	var conf configMap

	viper.SetConfigFile(filePath)
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	if err := viper.Unmarshal(&conf); err != nil {
		return nil, err
	}

	LoginTimeLimit, err := time.ParseDuration(conf.LoginTimeLimit)
	if err != nil {
		log.Fatalf("error reading configMap: %v", err)
	}

	IPTimeLimit, err := time.ParseDuration(conf.IPTimeLimit)
	if err != nil {
		log.Fatalf("error reading configMap: %v", err)
	}

	PasswordTimeLimit, err := time.ParseDuration(conf.PasswordTimeLimit)
	if err != nil {
		log.Fatalf("error reading configMap: %v", err)
	}

	CleanUpPeriod, err := time.ParseDuration(conf.CleanupPeriod)
	if err != nil {
		log.Fatalf("error reading configMap: %v", err)
	}

	return &Config{
		IPTimeLimit,
		conf.IPMaxCountForTimeLimit,
		LoginTimeLimit,
		conf.LoginMaxCountForTimeLimit,
		PasswordTimeLimit,
		conf.PasswordMaxCountForTimeLimit,
		conf.GrpcHost,
		conf.GrpcPort,
		CleanUpPeriod,
	}, nil
}
