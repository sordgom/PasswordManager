package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	ServerAddress string `mapstructure:"SERVER_ADDRESS"`

	Environment       string `mapstructure:"ENVIRONMENT"`
	RedisAddress      string `mapstructure:"REDIS_ADDRESS"`
	HTTPServerAddress string `mapstructure:"HTTP_SERVER_ADDRESS"`
	TokenSymmetricKey string `mapstructure:"TOKEN_SYMMETRIC_KEY"`
}

func LoadConfig(path string) (config Config, err error) {
	viper.SetConfigName("app")
	viper.SetConfigType("env")
	viper.AddConfigPath(path)

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	return
}
