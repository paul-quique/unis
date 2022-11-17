package api

import (
	"time"

	"github.com/jmoiron/sqlx"
)

const (
	OfferExpirationTime = 24 * time.Hour
	INSERT_OFFER        = "INSERT INTO offer (borrower_id, lender_id, product_id, created_at, expires_at) VALUES (:borrower_id, :lender_id, :product_id, :created_at, expires_at);"
)

type Offer struct {
	Id         int       `json:"id" db:"id"`
	BorrowerId int       `json:"borrowerId" db:"borrower_id"`
	LenderId   int       `json:"lenderId" db:"lender_id"`
	ProductId  int       `json:"productId" db:"product_id"`
	CreatedAt  time.Time `json:"createdAt" db:"created_at"`
	ExpiresAt  time.Time `json:"expiresAt" db:"expires_at"`
}

func (o *Offer) CreateInDB(db *sqlx.DB) (err error) {
	_, err = db.NamedExec(INSERT_OFFER, o)
	return err
}
