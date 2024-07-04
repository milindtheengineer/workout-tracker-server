package config

import (
	"fmt"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

var AppConfig Config

func InitialiseConfig() error {
	godotenv.Load("config.env")
	err := envconfig.Process("server", &AppConfig)
	if err != nil {
		return fmt.Errorf("InitializeConfig: %v", err)
	}
	return nil
}
