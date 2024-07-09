package database

type User struct {
	Email string
	Name  string
}

type UserRow struct {
	Id int
	User
}

type Session struct {
	UserID   int
	DateTime string
}

type SessionRow struct {
	Id int
	Session
}

type Workout struct {
	SessionID   int
	WorkoutName string
}

type WorkoutRow struct {
	Id int
	Workout
}

type Set struct {
	WorkoutID    int
	Weight       float32
	NumberOfReps int
}

type SetRow struct {
	Id int
	Set
}
