package config

import (
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	JWTPublicKey string `mapstructure:"JWT_PUBLIC_KEY"`
	PostgresUrl  string `mapstructure:"POSTGRES_URL"`
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

	JWTPublicKey, ok := viper.Get("JWT_PUBLIC_KEY").(string)
	if !ok {
		log.Fatalf("JWT_PUBLIC_KEY setting required")
	}

	postgresUrl, ok := viper.Get("POSTGRES_URL").(string)
	if !ok {
		log.Fatalf("POSTGRES_URL setting required")
	}

	return &Config{
		JWTPublicKey: JWTPublicKey,
		PostgresUrl:  postgresUrl,
	}, nil
}
