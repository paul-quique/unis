package api

import (
	"time"
)

type User struct {
	Id        int       `json:"id" db:"id"`
	FirstName string    `json:"firstName" db:"first_name"`
	LastName  string    `json:"lastName" db:"last_name"`
	Email     string    `json:"email" db:"email"`
	Password  string    `json:"password" db:"salted_hash"`
	Points    int       `json:"points" db:"points"`
	CreatedAt time.Time `json:"createdAt" db:"created_at"`
}
