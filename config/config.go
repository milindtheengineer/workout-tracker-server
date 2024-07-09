package config

import (
	"fmt"
	"runtime"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

var AppConfig Config

func InitialiseConfig() error {
	if runtime.GOOS != "linux" {
		godotenv.Load("config.env")
	}
	err := envconfig.Process("server", &AppConfig)
	if err != nil {
		return fmt.Errorf("InitializeConfig: %v", err)
	}
	return nil
}
