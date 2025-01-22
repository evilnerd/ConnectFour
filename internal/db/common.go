package db

import (
	"database/sql"
	log "github.com/sirupsen/logrus"
)

// readSecret reads a Docker secret from the designated location and returns the contents of the file as a string.
func readSecret(name string) string {

	// TODO: Implement reading the Docker secret or a local file
	return "dfjh4587fBBFgnjkfgno3mocopc0"
}

func connect() *sql.DB {
	var secret = readSecret("secret")
	db, err := sql.Open("mysql", "connectfour:"+secret+"@tcp(0.0.0.0:3306)/connectfour?parseTime=true")

	if err != nil {
		log.Fatalf("Error connecting to the database: %v\n", err)
	}
	return db
}
