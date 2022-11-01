package api

import "github.com/jmoiron/sqlx"

const (
	MaxProductNameLength = 45
	MinProductNameLength = 2
	INSERT_PRODUCT       = "INSERT INTO product (name, category_id, price) VALUES (:name, :category_id, :price);"
)

type Product struct {
	Id         int    `json:"id" db:"id"`
	Name       string `json:"name" db:"name"`
	CategoryId int    `json:"categoryId" db:"category_id"`
	Price      int    `json:"price" db:"price"`
}

func (p *Product) CreateInDB(db *sqlx.DB) error {
	_, err := db.NamedExec(INSERT_PRODUCT, p)
	return err
}
