package config

import (
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	Secret      string `mapstructure:"SECRET"`
	PostgresUrl string `mapstructure:"POSTGRES_URL"`
}

// TODO: create pydantic BaseSettings logic
func Read(tp, name, path string) (*Config, error) {
	// Environment variables
	viper.AutomaticEnv()

	// Configuration file
	viper.SetConfigType(tp)
	viper.SetConfigName(name)
	viper.AddConfigPath(path)

	// Read configuration
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, err
		}
	}

	secret, ok := viper.Get("SECRET").(string)
	if !ok {
		log.Fatalf("SECRET setting required")
	}

	postgresUrl, ok := viper.Get("POSTGRES_URL").(string)
	if !ok {
		log.Fatalf("POSTGRES_URL setting required")
	}

	return &Config{
		Secret:      secret,
		PostgresUrl: postgresUrl,
	}, nil
}
