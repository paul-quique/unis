package api

import (
	"database/sql"
	"fmt"
)

const (
	DbServer   = "tai.db.elephantsql.com"
	DbName     = "iuraljbb"
	DbUserName = "iuraljbb"
	DbPassword = "rncIyPl3pYQMTJlPQLDEiRgBP0BioWGR"
)

func ConnectToDatabase() error {
	connString := fmt.Sprintf("host=%s dbname=%s user=%s password=%s", DbServer, DbName, DbUserName, DbPassword)
	_, err := sql.Open("postgres", connString)
	if err != nil {
		return err
	}
	return nil
}
