package api

import (
	"fmt"
	"os"

	"github.com/jmoiron/sqlx"

	_ "github.com/lib/pq"
)

var (
	APIDatabase *sqlx.DB
)

func init() {
	var err error
	APIDatabase, err = ConnectToDatabase()
	if err != nil {
		panic(err)
	}
}

func ConnectToDatabase() (*sqlx.DB, error) {
	connString := fmt.Sprintf("host=%s dbname=%s user=%s password=%s",
		os.Getenv("DB_SERVER"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_USERNAME"),
		os.Getenv("DB_PASSWORD"))
	db, err := sqlx.Open("postgres", connString)
	if err != nil {
		return nil, err
	}
	return db, nil
}
