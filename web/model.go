package web

import "github.com/milindtheengineer/workout-tracker-server/database"

type Workout struct {
	database.WorkoutRow
	Sets []database.SetRow
}

type WorkoutIDResponse struct {
	WorkoutID int
}
