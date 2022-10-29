package api

import (
	"time"
)

type Offer struct {
	Id         int       `json:"id" db:"id"`
	BorrowerId int       `json:"borrowerId" db:"borrower_id"`
	LenderId   int       `json:"lenderId" db:"lender_id"`
	ProductId  int       `json:"productId" db:"product_id"`
	CreatedAt  time.Time `json:"createdAt" db:"created_at"`
	ExpiresAt  time.Time `json:"expiresAt" db:"expires_at"`
}
