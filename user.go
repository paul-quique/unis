package api

import (
	"net/http"
	"strconv"
	"time"

	"net/mail"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"golang.org/x/crypto/bcrypt"
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

type CreateUserRequest struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email" `
	Password  string `json:"password"`
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

func PostUser(c *gin.Context) {
	req := &CreateUserRequest{}
	if err := c.BindJSON(req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "please provide valid user informations",
		})
		return
	}

	//vérifier que le nom et le prénom ne sont pas nuls
	if req.FirstName == "" || req.LastName == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "please provide valid first name and last name",
		})
		return
	}

	//vérifier que l'email est de la forme d'une adresse mail
	if _, err := mail.ParseAddress(req.Email); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "please provide a valid email address",
		})
		return
	}

	//vérifier que le mot de passe est assez long
	if len(req.Password) < 8 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "password should be at least 8 characters long",
		})
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "internal server error occured, please try again later",
		})
		return
	}

	u := &User{
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Email:     req.Email,
		Password:  string(hash),
		CreatedAt: time.Now(),
		Points:    100,
	}

	err = u.CreateInDB(APIDatabase)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "user may already exist, please login or try with another email",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"createdAt": u.CreatedAt,
	})
}
