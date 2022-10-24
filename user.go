package api

import (
	"time"
)

type User struct {
	Id        int       `json:"id"`
	FirstName string    `json:"firstName"`
	LastName  string    `json:"lastName"`
	Email     string    `json:"email"`
	Password  string    `json:"password"`
	Points    int       `json:"points"`
	CreatedAt time.Time `json:"createdAt"`
}
