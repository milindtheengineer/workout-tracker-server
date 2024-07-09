package main

import (
	"github.com/milindtheengineer/workout-tracker-server/config"
	"github.com/milindtheengineer/workout-tracker-server/web"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	if err := config.InitialiseConfig(); err != nil {
		log.Panic().Msgf("Config could not be initialized due to %v", err)
	}
	log.Info().Msgf("Config is %v", config.AppConfig)
	web.StartRouter()
}
