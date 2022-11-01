package api

import "github.com/jmoiron/sqlx"

const (
	INSERT_CATEGORY = "INSERT INTO category (name) VALUES (:name);"
)

type Category struct {
	Id   int    `json:"id" db:"id"`
	Name string `json:"name" db:"name"`
}

func (c *Category) CreateInDb(db *sqlx.DB) error {
	_, err := db.NamedExec(INSERT_CATEGORY, c)
	return err
}
