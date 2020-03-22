// https://github.com/mattn/go-sqlite3/issues/274
// https://github.com/mattn/go-sqlite3/issues/569

package main

import (
	"database/sql"
	"fmt"
	"strconv"
	"strings"

	"bufio"
	"os"

	"flag"

	_ "github.com/mattn/go-sqlite3"
	"steadylearner.com/sqlite/models"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"steadylearner.com/sqlite/repo"
)

var (
	action string
)

func important(err error, reason string) {
	if err != nil {
		log.Fatal(errors.Wrap(err, reason))
		return
	}
}

func init() {
	actionHelp := "Use [list, get, create, update, delete] to handle the data from users."
	flag.StringVar(&action, "action", "list", actionHelp)
	flag.Parse()
}

func main() {
	target := "database/users.db"

	db, err := sql.Open("sqlite3", target)
	important(err, fmt.Sprintf("Couldn't open %s", target))
	defer db.Close()

	// https://www.sqlitetutorial.net/sqlite-autoincrement/
	db.Exec("CREATE TABLE IF NOT EXISTS users (id INTEGER PRIMARY KEY AUTOINCREMENT, name TEXT NOT NULL UNIQUE)")

	// db.Exec("DROP TABLE users")
	// db.Exec("CREATE TABLE users (id INTEGER PRIMARY KEY AUTOINCREMENT, username TEXT NOT NULL UNIQUE)")

	userRepo := repo.NewUsersRepository(db)

	switch action {
	case "list":
		list, err := userRepo.List()
		important(err, "fail to read user list")

		for _, u := range list {
			fmt.Printf("User [%d]: %s\n", u.ID, u.Name)
		}
	case "get":
		idReader := bufio.NewReader(os.Stdin)
		fmt.Println("Type the id of a user get the data.")
		id, err := idReader.ReadString('\n')
		important(err, "Fail to read the id to get the data of a user.") // Substitute it with logrus

		if id == "\n" {
			fmt.Println("You should provide an id to make this work.")
			return
		}

		i, err := strconv.Atoi(strings.TrimSuffix(id, "\n"))
		important(err, fmt.Sprintf("Fail to convert given id(%s) to int.", id)) // Substitute it with logrus

		user, err := userRepo.Get(int64(i))
		important(err, fmt.Sprintf("Fail to feetch user by ID (%s) to int.", id)) // Substitute it with logrus

		fmt.Printf("%+v\n", user)
	case "create":
		usernameReader := bufio.NewReader(os.Stdin)
		fmt.Println("Type a username you want to use.")
		username, err := usernameReader.ReadString('\n')
		important(err, "Fail to read a username.") // Substitute it with logrus

		if username == "\n" {
			fmt.Println("You should provide username to make it work.")
			return
		}

		newUser := models.User{Name: strings.TrimSpace(username)}
		err = userRepo.Create(&newUser)
		important(err, fmt.Sprintf("Fail to create user(%s).", username)) // Substitute it with logrus

		fmt.Printf("User %d created.\n", newUser.ID)
	case "update":
		idReader := bufio.NewReader(os.Stdin)
		fmt.Println("Type the id of a user to update its username.")
		id, err := idReader.ReadString('\n')
		important(err, "Fail to read the id of a user to be updated.") // Substitute it with logrus

		if id == "\n" {
			fmt.Println("You should provide an id to make this work.")
			return
		}

		newUsernameReader := bufio.NewReader(os.Stdin)
		fmt.Println("Type a new username you want to use.")
		newUsername, err := newUsernameReader.ReadString('\n')
		important(err, "Fail to read a new username.") // Substitute it with logrus

		if newUsername == "\n" {
			fmt.Println("You should provide username to make it work.")
			return
		}

		i, err := strconv.Atoi(strings.TrimSuffix(id, "\n"))
		important(err, fmt.Sprintf("Fail to convert given id(%s) to int.", id)) // Substitute it with logrus

		upd := models.User{ID: int64(i), Name: strings.TrimSpace(newUsername)}
		err = userRepo.Update(&upd)
		important(err, "Fail to update user.") // Substitute it with logrus

	case "delete":
		idReader := bufio.NewReader(os.Stdin)
		fmt.Println("Type the id of a user to delete.")
		id, err := idReader.ReadString('\n')
		important(err, "Fail to read the id to delete the data of a user.") // Substitute it with logrus

		if id == "\n" {
			fmt.Println("You should provide the id to make it work.")
			return

		}
		i, err := strconv.Atoi(strings.TrimSuffix(id, "\n"))
		important(err, fmt.Sprintf("Fail to convert given id(%s) to int.", id)) // Substitute it with logrus

		err = userRepo.Delete(int64(i))
		important(err, fmt.Sprintf("Fail to delete user (%d)", i))
	default:
		fmt.Println("You should use [list, get, create, update, delete].")
	}
}
