package api

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

const (
	GET_USER_BY_ID    = "SELECT * FROM user_info WHERE id=$1;"
	GET_USER_BY_EMAIL = "SELECT * FROM user_info WHERE email=$1;"
	INSERT_USER       = "INSERT INTO user_info (first_name, last_name, email, salted_hash, points, created_at) VALUES (:first_name, :last_name, :email, :salted_hash, :points, :created_at);"
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

func (u *User) CreateInDB(db *sqlx.DB) error {
	_, err := db.NamedExec(INSERT_USER, u)
	return err
}

func LoadUserFromId(db *sqlx.DB, id int) (*User, error) {
	u := &User{}
	err := db.Get(u, GET_USER_BY_ID, id)
	if err != nil {
		return nil, err
	}
	return u, nil
}

func LoadUserFromEmail(db *sqlx.DB, email string) (*User, error) {
	u := &User{}
	err := db.Get(u, GET_USER_BY_EMAIL, email)
	if err != nil {
		return nil, err
	}
	return u, nil
}

func GetUser(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "user id should be an integer",
		})
		return
	}

	u, err := LoadUserFromId(APIDatabase, id)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "no user matches the given id",
		})
		return
	}

	c.JSON(http.StatusFound, gin.H{
		"id":        u.Id,
		"firstName": u.FirstName,
		"lastName":  u.LastName,
		"email":     u.Email,
		"points":    u.Points,
		"createdAt": u.CreatedAt,
	})
}
