package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
)

type Store struct {
	DB *sql.DB
}

var DB *sql.DB

var (
	host     = "localhost"
	port     = 5432
	password = ""
	dbname   = "fileserver"
)

func Conn() error {
	password = os.Getenv("PGPASSWORD")
	connOptions := fmt.Sprintf("host=%s user=postgres port=%d password=%s sslmode=disable", host, port, password)

	db, err := sql.Open("postgres", connOptions)
	if err != nil {
		log.Println("Could not open db")
		return err
	}

	err = db.Ping()
	if err != nil {
		log.Println("Could not ping db", err)
		return err
	}

	log.Println("Connected")
	DB = db
	return nil
}
