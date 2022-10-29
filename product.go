package api

import "fmt"

const (
	MaxProductNameLength = 45
	MinProductNameLength = 2
)

type Product struct {
	Id         int    `json:"id" db:"id"`
	Name       string `json:"name" db:"name"`
	CategoryId int    `json:"categoryId" db:"category_id"`
	Price      int    `json:"price" db:"price"`
}

func Validate(p *Product) error {
	if p.Price <= 0 {
		return fmt.Errorf("the price must be positive, price: %d", p.Price)
	}
	l := len([]rune(p.Name))
	if l > MaxProductNameLength {
		return fmt.Errorf("the name shouldn't be longer than %d, name length: %d", MaxProductNameLength, l)
	}
	if l < MinProductNameLength {
		return fmt.Errorf("the name shouldn't be shorter than %d, name length: %d", MinProductNameLength, l)
	}
	return nil
}
