package web

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/milindtheengineer/workout-tracker-server/config"
	"github.com/milindtheengineer/workout-tracker-server/database"
	"github.com/rs/zerolog"
)

func HealthHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Ok"))
}

type App struct {
	db     *database.DBConn
	logger zerolog.Logger
}

func authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Authorization header required", http.StatusUnauthorized)
			return
		}
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString != config.AppConfig.Token {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}

// All gets

// Get Sessions based on userID (restrict to 10 in the future maybe)
func (app *App) SessionHandler(w http.ResponseWriter, r *http.Request) {
	userIDStr := chi.URLParam(r, "userID")
	if len(userIDStr) < 1 {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}
	sessions, err := app.db.GetSessionsByUserId(userID)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		app.logger.Error().Msgf("SessionHandler: %v", err)
		return
	}
	body, err := json.Marshal(sessions)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		app.logger.Error().Msgf("SessionHandler: %v", err)
		return
	}
	w.Write(body)
}

// Get Workouts based on sessionID (restrict to 10 in the future maybe)
func (app *App) WorkoutHandler(w http.ResponseWriter, r *http.Request) {
	var workoutResponse []Workout
	sessionIDstr := chi.URLParam(r, "sessionID")
	if len(sessionIDstr) < 1 {
		http.Error(w, "Invalid session ID", http.StatusBadRequest)
		return
	}
	sessionID, err := strconv.Atoi(sessionIDstr)
	if err != nil {
		http.Error(w, "Invalid session ID", http.StatusBadRequest)
		return
	}
	workouts, err := app.db.GetWorkoutsBySessionId(sessionID)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		app.logger.Error().Msgf("WorkoutHandler: %v", err)
		return
	}
	for _, workout := range workouts {
		sets, err := app.db.GetSetsByWorkoutId(workout.Id)
		if err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			app.logger.Error().Msgf("WorkoutHandler: %v", err)
			return
		}
		workoutResponse = append(workoutResponse, Workout{WorkoutRow: workout, Sets: sets})
	}
	body, err := json.Marshal(workoutResponse)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		app.logger.Error().Msgf("WorkoutHandler: %v", err)
		return
	}
	w.Write(body)
}

// Get Workouts and sets based on sessionID
