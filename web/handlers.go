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
func (app *App) SessionListHandler(w http.ResponseWriter, r *http.Request) {
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

// Get Workouts based on sessionId
func (app *App) WorkoutListHandler(w http.ResponseWriter, r *http.Request) {
	workoutResponse := []Workout{}
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
		if len(sets) == 0 {
			sets = []database.SetRow{}
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

// Get sets based on workoutID
func (app *App) SetListHandler(w http.ResponseWriter, r *http.Request) {
	workoutIDstr := chi.URLParam(r, "workoutID")
	if len(workoutIDstr) < 1 {
		http.Error(w, "Invalid workout ID", http.StatusBadRequest)
		return
	}
	workoutID, err := strconv.Atoi(workoutIDstr)
	if err != nil {
		http.Error(w, "Invalid workout ID", http.StatusBadRequest)
		return
	}
	sets, err := app.db.GetSetsByWorkoutId(workoutID)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		app.logger.Error().Msgf("SetListHandler: %v", err)
		return
	}
	body, err := json.Marshal(sets)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		app.logger.Error().Msgf("SetListHandler: %v", err)
		return
	}
	w.Write(body)
}

func (app *App) WorkoutCreateHandler(w http.ResponseWriter, r *http.Request) {
	var workout database.Workout
	if err := json.NewDecoder(r.Body).Decode(&workout); err != nil {
		http.Error(w, "Could not decode workout", http.StatusBadRequest)
		return

	}
	if err := app.db.CreateWorkoutForSession(workout.SessionID, strings.ToLower(workout.WorkoutName), workout.UserID); err != nil {
		http.Error(w, "Could not decode workout", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (app *App) SetCreateHandler(w http.ResponseWriter, r *http.Request) {
	var set database.Set
	if err := json.NewDecoder(r.Body).Decode(&set); err != nil {
		app.logger.Error().Msgf("%v", err)
		http.Error(w, "Could not decode set", http.StatusBadRequest)
		return

	}
	if err := app.db.CreateSetForWorkout(set.WorkoutID, set.NumberOfReps, set.Weight); err != nil {
		app.logger.Error().Msgf("%v", err)
		http.Error(w, "Could not add set to workout", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (app *App) SessionCreateHandler(w http.ResponseWriter, r *http.Request) {
	var session database.Session
	if err := json.NewDecoder(r.Body).Decode(&session); err != nil {
		app.logger.Error().Msgf("%v", err)
		http.Error(w, "Could not decode session", http.StatusBadRequest)
		return

	}
	if err := app.db.CreateSessionForUser(session.UserID); err != nil {
		app.logger.Error().Msgf("%v", err)
		http.Error(w, "Could not add session", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

// Get Sessions based on userID (restrict to 10 in the future maybe)
func (app *App) LastWorkoutHandler(w http.ResponseWriter, r *http.Request) {
	userIDStr := chi.URLParam(r, "userID")
	if len(userIDStr) < 1 {
		http.Error(w, "Invalid user id", http.StatusBadRequest)
		return
	}
	workoutName := chi.URLParam(r, "workout")
	if len(workoutName) < 1 {
		http.Error(w, "Invalid workout name", http.StatusBadRequest)
		return
	}
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}
	workoutID, err := app.db.GetLastWorkoutID(workoutName, userID)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		app.logger.Error().Msgf("LastWorkoutHandler: %v", err)
		return
	}
	var lastWorkoutDetails []database.SetRow
	if workoutID > 0 {
		lastWorkoutDetails, err = app.db.GetSetsByWorkoutId(workoutID)
		if err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			app.logger.Error().Msgf("LastWorkoutHandler: %v", err)
			return
		}
	}
	body, err := json.Marshal(lastWorkoutDetails)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		app.logger.Error().Msgf("LastWorkoutHandler: %v", err)
		return
	}
	w.Write(body)
}
