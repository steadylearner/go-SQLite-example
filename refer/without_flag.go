// https://github.com/mattn/go-sqlite3/issues/274
// https://github.com/mattn/go-sqlite3/issues/569

package main

import (
	"database/sql"
	"fmt"
	"strconv"

	_ "github.com/mattn/go-sqlite3"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	Models "steadylearner.com/sqlite/models"
)

func important(err error, reason string) {
	if err != nil {
		log.Fatal(errors.Wrap(err, reason))
		return
	}
}

func createUser(db *sql.DB, username string) {
	tx, err := db.Begin()
	important(err, fmt.Sprintf("Fail to begin db to create the user %s.", username)) // Substitute it with logrus.

	cmd := "INSERT INTO users (username) values (?)"
	stmt, err := tx.Prepare(cmd)
	important(err, fmt.Sprintf("Fail to prepare sql statement(%s).", cmd)) // Substitute it with logrus.
	defer stmt.Close()

	// https://godoc.org/database/sql/driver#Result
	_, err = stmt.Exec(username)                                        // What is ok value here?
	important(err, fmt.Sprintf("Fail to create the user %s", username)) // Substitute it with logrus.

	// https://github.com/mattn/go-sqlite3/blob/master/sqlite3.go#L443
	tx.Commit()
}

func updateUser(db *sql.DB, username string, id int) {
	sid := strconv.Itoa(id) // int to string

	tx, err := db.Begin()
	important(err, fmt.Sprintf("Fail to begin db to update the username to %s.", username)) // Substitute it with logrus.

	cmd := "UPDATE users SET username=? where id=?"
	stmt, err := db.Prepare(cmd)
	important(err, fmt.Sprintf("Fail to prepare sql statement(%s).", cmd)) // Substitute it with logrus.
	defer stmt.Close()

	result, err := stmt.Exec(username, sid)
	important(err, fmt.Sprintf("Fail to update the user %s.", username)) // Substitute it with logrus.

	updatedRows, err := result.RowsAffected()
	important(err, fmt.Sprintf("Fail to find the number of rows affected with %s. ", cmd)) // Substitute it with logrus.

	fmt.Println(updatedRows)

	tx.Commit()
}

func getUser(db *sql.DB, id int) Models.User { // Models.User
	sid := strconv.Itoa(id) // int to string

	tx, err := db.Begin()
	important(err, fmt.Sprintf("Fail to begin db to get the user with id %s.", sid)) // Substitute it with logrus.

	cmd := fmt.Sprintf("SELECT * FROM users WHERE id=%s", sid)
	rows, err := tx.Query(cmd)                                           // https://godoc.org/database/sql/driver#Rows
	important(err, fmt.Sprintf("Fail to query sql statement(%s).", cmd)) // Substitute it with logrus.
	defer rows.Close()

	if rows.Next() {
		var user Models.User

		err = rows.Scan(&user.Id, &user.Username)
		important(err, fmt.Sprintf("Fail to convert data from SQLite to the Models.User type.")) // Substitute it with logrus.
		return user
	}
	err = rows.Err()
	important(err, "Problem while reading a row in getUser function.")

	return Models.User{}
}

func listUsers(db *sql.DB) []Models.User {
	rows, err := db.Query("SELECT * FROM users")
	important(err, "Fail to list users") // Substitute it with logrus.

	var users []Models.User
	for rows.Next() {
		var user Models.User
		err = rows.Scan(&user.Id, &user.Username)
		important(err, fmt.Sprintf("Fail to convert data from SQLite to the Models.User type.")) // Substitute it with logrus.
		users = append(users, user)
	}
	err = rows.Err()
	important(err, "Problem while reading rows in listUsers function.")

	return users
}

func deleteUser(db *sql.DB, id int) {
	tx, err := db.Begin()
	important(err, fmt.Sprintf("Fail to begin db to delete the user with id %d.", id)) // Substitute it with logrus.

	sid := strconv.Itoa(id)

	cmd := "DELETE FROM users WHERE id=?"

	stmt, err := tx.Prepare(cmd)
	important(err, fmt.Sprintf("Fail to prepare sql statement(%s).", cmd)) // Substitute it with logrus.
	result, err := stmt.Exec(sid)
	important(err, fmt.Sprintf("Fail to delete the user with id %d", id)) // Substitute it with logrus
	num, err := result.RowsAffected()
	important(err, fmt.Sprintf("Fail to find the number of rows affected with %s", cmd))
	fmt.Printf("%d user deleted.\n", num)

	tx.Commit()
}

// 2. Make CLI.
// 3. Include a log file and relevant code.
// 4. Create a web app with more fields.

// func init() {
// 	// Log as JSON instead of the default ASCII formatter.
// 	// log.SetFormatter(&log.JSONFormatter{})
// 	log.SetFormatter(&log.TextFormatter{
// 		// DisableColors: true,
// 		FullTimestamp: true,
// 	})
// }

// defaultUserLogger := log.WithFields(log.Fields{
// 			"username": m.Sender.Username,
// 			"text":     m.Text,
// 		})

// defaultUserLogger.Info("Models.User request without /")
// 			log.SetOutput(f)

// defaultUserLogger.WithError(err).Info(info)

// f, err := os.OpenFile("telebot.log",
// 		os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
// 	important(err, "You should verify the telebot.log file.")
// 	defer f.Close()

// Verify the result with $sqlite3 users.db (https://sqlite.org/cli.html)
func main() {
	// $mkdir database or write code for it.
	target := "database/users.db"
	// os.Remove(target)

	db, err := sql.Open("sqlite3", target)
	defer db.Close()

	important(err, fmt.Sprintf("Couldn't open %s", target))

	// https://www.sqlitetutorial.net/sqlite-autoincrement/
	// db.Exec("CREATE TABLE IF NOT EXISTS users (id INTEGER PRIMARY KEY AUTOINCREMENT, username TEXT NOT NULL UNIQUE)")

	// db.Exec("DROP TABLE users")
	// db.Exec("CREATE TABLE users (id INTEGER PRIMARY KEY AUTOINCREMENT, username TEXT NOT NULL UNIQUE)")

	// createUser(db, "Steady")
	// createUser(db, "Learner")

	// updateUser(db, "Gopher", 1)
	// updateUser(db, "Pythonist", 2)

	// fmt.Println(getUser(db, 1))
	// fmt.Println(getUser(db, 2))

	// deleteUser(db, 1)
	// defer fmt.Println(getUser(db, 1))

	// fmt.Println(listUsers(db))
}
