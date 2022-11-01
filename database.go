package api

import (
	"fmt"

	"github.com/jmoiron/sqlx"

	_ "github.com/lib/pq"
)

var (
	APIDatabase *sqlx.DB
)

const (
	//Ces données seront des variables d'environnement
	DbServer   = "tai.db.elephantsql.com"
	DbName     = "iuraljbb"
	DbUserName = "iuraljbb"
	DbPassword = "rncIyPl3pYQMTJlPQLDEiRgBP0BioWGR"
)

func init() {
	var err error
	APIDatabase, err = ConnectToDatabase()
	if err != nil {
		panic(err)
	}
}

func ConnectToDatabase() (*sqlx.DB, error) {
	connString := fmt.Sprintf("host=%s dbname=%s user=%s password=%s", DbServer, DbName, DbUserName, DbPassword)
	db, err := sqlx.Open("postgres", connString)
	if err != nil {
		return nil, err
	}
	return db, nil
}
