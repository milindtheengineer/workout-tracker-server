package database

import (
	"database/sql"
	"fmt"
	"time"

	"modernc.org/sqlite"
	sqlite3 "modernc.org/sqlite/lib"
)

type DBConn struct {
	db *sql.DB
}

func CreateDBConnection(dbPath string) (*DBConn, error) {
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, fmt.Errorf("CreateDBConnection: %w", err)
	}
	return &DBConn{db: db}, err
}

func (d *DBConn) CloseConn() error {
	return d.db.Close()
}

func (d *DBConn) CreateUser(user User) (int64, error) {
	stmt, err := d.db.Prepare("INSERT INTO User (email, name) VALUES (?, ?)")
	if err != nil {
		return 0, fmt.Errorf("CreateUser: error preparing statement: %w", err)
	}
	defer stmt.Close()

	result, err := stmt.Exec(user.Email, user.Name)
	if err != nil {
		return 0, fmt.Errorf("CreateUser: error executing statement: %w", err)
	}
	userID, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("CreateUser: Error getting last insert ID: %w", err)
	}
	return userID, nil
}

func (d *DBConn) GetUserByEmail(email string) (UserRow, error) {
	// Query to get a user by email
	query := "SELECT userId, email, name FROM User WHERE email = ?"
	var user UserRow

	// Execute the query with the specified email
	if err := d.db.QueryRow(query, email).Scan(&user.Id, &user.Email, &user.Name); err != nil {
		if err == sql.ErrNoRows {
			return user, fmt.Errorf("no user found") // TODO: do a not found error later
		}
		return user, fmt.Errorf("GetUserByEmail: %v", err)
	}
	return user, nil
}

func (d *DBConn) CreateSessionForUser(userID int) error {
	// Prepare the insert statement
	stmt, err := d.db.Prepare("INSERT INTO Session (userID, dateTime) VALUES (?, ?)")
	if err != nil {
		return fmt.Errorf("CreateSessionForUser: %v", err)
	}
	defer stmt.Close()

	dateTime := time.Now().Format("2006-01-02 15:04:05")

	// Execute the insert statement
	_, err = stmt.Exec(userID, dateTime)
	if err != nil {
		return fmt.Errorf("CreateSessionForUser: %w", err)
	}
	return nil
}

func (d *DBConn) GetSessionsByUserId(userId int) ([]SessionRow, error) {
	// Query to get all sessions for the specified userID
	query := "SELECT sessionID, userID, dateTime FROM Session WHERE userID = ? ORDER BY sessionID DESC"

	// Execute the query
	rows, err := d.db.Query(query, userId)
	if err != nil {
		return nil, fmt.Errorf("GetSessionsByUserId: Error executing query: %w", err)
	}
	defer rows.Close()

	// Iterate over the result set
	var sessions []SessionRow
	for rows.Next() {
		var sess SessionRow
		err := rows.Scan(&sess.Id, &sess.UserID, &sess.DateTime)
		if err != nil {
			return nil, fmt.Errorf("GetSessionsByUserId: Error scanning row: %w", err)
		}
		sessions = append(sessions, sess)
	}

	// Check for errors after iterating over rows
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("GetSessionsByUserId: Error after iterating rows: %w", err)
	}
	return sessions, nil
}

func (d *DBConn) GetWorkoutsBySessionId(sessionId int) ([]WorkoutRow, error) {
	// Query to get workouts for the specified sessionID
	query := "SELECT workoutID, sessionID, workoutname FROM Workouts WHERE sessionID = ? ORDER BY workoutID DESC"

	// Execute the query
	rows, err := d.db.Query(query, sessionId)
	if err != nil {
		return nil, fmt.Errorf("GetWorkoutsBySessionId: Error executing query: %w", err)
	}
	defer rows.Close()

	// Iterate over the result set
	var workouts []WorkoutRow
	for rows.Next() {
		var workout WorkoutRow
		err := rows.Scan(&workout.Id, &workout.SessionID, &workout.WorkoutName)
		if err != nil {
			return nil, fmt.Errorf("Error scanning row: %w", err)
		}
		workouts = append(workouts, workout)
	}

	// Check for errors after iterating over rows
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("Error after iterating rows: %w", err)
	}

	return workouts, nil
}

func (d *DBConn) CreateWorkoutForSession(sessionId int, workoutname string, userId int) error {
	stmt, err := d.db.Prepare("INSERT INTO Workouts (sessionID, workoutname, userID) VALUES (?, ?, ?)")
	if err != nil {
		return fmt.Errorf("Error preparing statement: %w", err)
	}
	defer stmt.Close()

	// Execute the insert statement
	_, err = stmt.Exec(sessionId, workoutname, userId)
	if err != nil {
		sqliteErr, ok := err.(*sqlite.Error)
		if ok && sqliteErr.Code() == sqlite3.SQLITE_CONSTRAINT_UNIQUE {
			return fmt.Errorf("Workout already exists: workout with name %s for session %d already exists for user %d", workoutname, sessionId, userId)
		}

		return fmt.Errorf("Error inserting new workout: %w", err)
	}
	return nil
	// // Get the last inserted ID (workoutID)
	// workoutID, err := result.LastInsertId()
	// if err != nil {
	//     fmt.Println("Error getting last insert ID:", err)
	//     return
	// }

	// fmt.Printf("New workout created successfully with workoutID: %d\n", workoutID)

}

func (d *DBConn) CreateSetForWorkout(workoutId int, numberofReps int, weight float32) error {
	stmt, err := d.db.Prepare("INSERT INTO Sets (numberofReps, weight, workoutID) VALUES (?, ?, ?)")
	if err != nil {
		return fmt.Errorf("Error preparing statement: %w", err)
	}
	defer stmt.Close()

	// Execute the insert statement
	_, err = stmt.Exec(numberofReps, weight, workoutId)
	if err != nil {
		return fmt.Errorf("Error inserting new set: %w", err)
	}

	// // Get the last inserted ID (setID)
	// setID, err := result.LastInsertId()
	// if err != nil {
	//     fmt.Println("Error getting last insert ID:", err)
	//     return
	// }

	// fmt.Printf("New set created successfully with setID: %d\n", setID)
	return nil
}

func (d *DBConn) GetSetsByWorkoutId(workoutID int) ([]SetRow, error) {
	// Query to get sets for the specified workoutID
	query := "SELECT setID, numberofReps, weight, workoutID FROM Sets WHERE workoutID = ? ORDER BY setID DESC"

	// Execute the query
	rows, err := d.db.Query(query, workoutID)
	if err != nil {
		return nil, fmt.Errorf("GetSetsByWorkoutId: %w", err)
	}
	defer rows.Close()

	// Iterate over the result set
	var sets []SetRow
	for rows.Next() {
		var set SetRow
		err := rows.Scan(&set.Id, &set.NumberOfReps, &set.Weight, &set.WorkoutID)
		if err != nil {
			return nil, fmt.Errorf("GetSetsByWorkoutId: %w", err)
		}
		sets = append(sets, set)
	}

	// Check for errors after iterating over rows
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("GetSetsByWorkoutId: %w", err)
	}
	return sets, nil
}

func (d *DBConn) GetLastWorkoutID(workoutName string, userID int) (int, error) {
	var workoutID int
	query := `
        SELECT workoutID
        FROM Workouts
        WHERE workoutName = ? AND userID = ?
        ORDER BY workoutID DESC
        LIMIT 1 OFFSET 1
    `
	if err := d.db.QueryRow(query, workoutName, userID).Scan(&workoutID); err != nil && err != sql.ErrNoRows {
		return workoutID, err
	}
	return workoutID, nil
}
