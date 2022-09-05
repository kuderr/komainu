package config

import (
	_ "embed"

	"github.com/spf13/viper"
)

type Config struct {
	Secret      string `mapstructure:"SECRET"`
	PostgresUrl string `mapstructure:"POSTGRES_URL"`
}

func Read() (*Config, error) {
	// Environment variables
	viper.AutomaticEnv()

	// Configuration file
	viper.SetConfigType("env")
	viper.SetConfigName(".env")
	viper.AddConfigPath(".")

	// Read configuration
	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	// Unmarshal the configuration
	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, err
	}
	return &config, nil
}
