package web

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/milindtheengineer/workout-tracker-server/database"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func StartRouter() {
	r := chi.NewRouter()
	r.Use(authMiddleware)
	r.Use(middleware.Logger)
	// Shift this logic to main probably
	db, err := database.CreateDBConnection("/Users/milindjuttiga/code/sqlite3db/test.db")
	if err != nil {
		panic(err)
	}
	app := App{
		db:     db,
		logger: zerolog.Logger{},
	}
	r.Get("/health", HealthHandler)
	r.Get("/sessions/{userID}", app.SessionHandler)
	r.Get("/workouts/{sessionID}", app.WorkoutHandler)
	// r.GET("/v1/user", authMiddleware(user.Crud))
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Panic().Msg(err.Error())
	}
}
