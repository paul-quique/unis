package api

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/lib/pq"
)

const (
	DbServer   = "tai.db.elephantsql.com"
	DbName     = "iuraljbb"
	DbUserName = "iuraljbb"
	DbPassword = "rncIyPl3pYQMTJlPQLDEiRgBP0BioWGR"
)

func ConnectToDatabase() error {
	connString := fmt.Sprintf("host=%s dbname=%s user=%s password=%s", DbServer, DbName, DbUserName, DbPassword)
	db, err := sql.Open("postgres", connString)
	if err != nil {
		return err
	}
	c, err := os.ReadFile("./sql/create_db.sql")
	if err != nil {
		return err
	}

	db.Exec(string(c))
	return nil
}
