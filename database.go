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
	host, dbname, user, password := os.Getenv("DB_SERVER"), os.Getenv("DB_NAME"), os.Getenv("DB_USERNAME"), os.Getenv("DB_PASSWORD")
	connString := fmt.Sprintf("host=%s dbname=%s user=%s password=%s",
		host,
		dbname,
		user,
		password)
	fmt.Println(host, dbname, user, password)
	db, err := sqlx.Open("postgres", connString)
	if err != nil {
		return nil, err
	}
	return db, nil
}
