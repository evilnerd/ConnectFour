package server

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	log "github.com/sirupsen/logrus"
)

var (
	db *sql.DB
)

// readSecret reads a Docker secret from the designated location and returns the contents of the file as a string.
func readSecret(name string) string {

	return "dfjh4587fBBFgnjkfgno3mocopc0"
}

func init() {
	var err error
	var secret = readSecret("secret")
	db, err = sql.Open("mysql", "connectfour:"+secret+"@tcp(0.0.0.0:3306)/connectfour")

	if err != nil {
		log.Fatalf("Error connecting to the database: %v\n", err)
	}
}

func CreateUser(u User) User {
	result, err := db.Exec("INSERT INTO user (email, name, token) VALUES (?, ?, ?)", u.Email, u.Name, u.Token)
	if err != nil {
		log.Errorf("Error inserting user into the database: %v\n", err)
		return User{}
	}
	id, err := result.LastInsertId()
	if err != nil {
		log.Errorf("Error getting last insert ID: %v\n", err)
		return User{}
	}
	u.Id = id
	return u
}
