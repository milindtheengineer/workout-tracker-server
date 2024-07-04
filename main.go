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
	// db, err := database.CreateDBConnection("/Users/milindjuttiga/code/sqlite3db/test.db")
	// if err != nil {
	// 	log.Panic().Msgf("Database creation error: %v", err)
	// }
	// defer db.CloseConn()
	// user, err := db.GetUserByEmail("milindjuttiga@gmail.com")
	// if err != nil {
	// 	log.Panic().Msgf("Database creation error: %v", err)
	// }
	// log.Info().Msgf("User %v", user)
	// if err := db.CreateSessionForUser(user.Id); err != nil {
	// 	log.Panic().Msgf("Create session error: %v", err)
	// }
	// sessions, err := db.GetSessionsByUserId(user.Id)
	// if err != nil {
	// 	log.Panic().Msgf("get session: %v", err)
	// }
	// if err := db.CreateWorkoutForSession(sessions[0].Id, "arm curl"); err != nil {
	// 	log.Panic().Msgf("get session: %v", err)
	// }
	// workouts, err := db.GetWorkoutsBySessionId(sessions[0].Id)
	// if err := db.CreateSetForWorkout(workouts[0].Id, 10, 30); err != nil {
	// 	log.Panic().Msgf("get session: %v", err)
	// }
	// sets, err := db.GetSetsByWorkoutId(workouts[0].Id)
	// for _, set := range sets {
	// 	fmt.Println(set)
	// }
	web.StartRouter()
}
