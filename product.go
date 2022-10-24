package api

import "fmt"

const (
	MaxProductNameLength = 45
	MinProductNameLength = 2
)

type Product struct {
	Id         int    `json:"id"`
	Name       string `json:"name"`
	CategoryId int    `json:"categoryId"`
	Price      int    `json:"price"`
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
