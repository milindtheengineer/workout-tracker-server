package web

import (
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/milindtheengineer/workout-tracker-server/config"
	"github.com/milindtheengineer/workout-tracker-server/database"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func StartRouter() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(cors.Handler(cors.Options{
		// AllowedOrigins:   []string{"https://foo.com"}, // Use this to allow specific origin hosts
		AllowedOrigins: []string{"https://workout-tracker.13059596.xyz"},
		// AllowOriginFunc:  func(r *http.Request, origin string) bool { return true },
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"X-PINGOTHER", "Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))
	// r.Use(authMiddleware)

	// Shift this logic to main probably
	db, err := database.CreateDBConnection(config.AppConfig.DBPath)
	if err != nil {
		panic(err)
	}
	app := App{
		db:     db,
		logger: zerolog.New(os.Stdout).With().Timestamp().Logger(),
	}
	r.Get("/health", HealthHandler)
	r.Post("/login", app.HandleLogin)
	r.Group(func(r chi.Router) {
		r.Use(app.authMiddleware)
		r.Get("/sessions", app.SessionListHandler)
		r.Get("/workouts/{sessionID}", app.WorkoutListHandler)
		r.Get("/sets/{workoutID}", app.SetListHandler)
		r.Get("/lastworkout/{workout}", app.LastWorkoutHandler)
		r.Post("/sessions", app.SessionCreateHandler)
		r.Post("/workouts", app.WorkoutCreateHandler)
		r.Post("/sets", app.SetCreateHandler)
	})

	// r.GET("/v1/user", authMiddleware(user.Crud))
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Panic().Msg(err.Error())
	}
}
