package db

import (
	"database/sql"
	log "github.com/sirupsen/logrus"
	"os"
)

// readSecret reads a Docker secret from the designated location and returns the contents of the file as a string.
func readSecret(name string) string {
	file := os.Getenv(name)
	b, err := os.ReadFile(file)
	if err != nil {
		log.Fatalf("Reading secret: %v\n", err)
	}
	return string(b)
}

func connect() *sql.DB {
	log.Infoln("Connecting to DB...")
	secret := readSecret("MARIADB_PASSWORD_FILE")
	user := os.Getenv("MARIADB_USER")
	schema := os.Getenv("MARIADB_DATABASE")
	address := os.Getenv("MARIADB_ADDRESS")
	datasource := user + ":" + secret + "@" + address + "/" + schema + "?parseTime=true"
	db, err := sql.Open("mysql", datasource)

	if err != nil {
		log.Fatalf("Error connecting to the database: %v\n", err)
	}
	log.Infoln("Connected.")
	return db
}
