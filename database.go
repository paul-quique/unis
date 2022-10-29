package api

import (
	"fmt"

	"github.com/jmoiron/sqlx"

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
	db, err := sqlx.Open("postgres", connString)
	if err != nil {
		return err
	}
	users := []User{}
	if err != nil {
		return err
	}
	db.Select(&users, "SELECT * FROM user_info;")
	fmt.Println(users)
	return nil
}
