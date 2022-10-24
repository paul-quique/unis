package api

import (
	"time"
)

type Offer struct {
	Id         int       `json:"id"`
	BorrowerId int       `json:"borrowerId"`
	LenderId   int       `json:"lenderId"`
	ProductId  int       `json:"productId"`
	CreatedAt  time.Time `json:"createdAt"`
	ExpiresAt  time.Time `json:"expiresAt"`
}
