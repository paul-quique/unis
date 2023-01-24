package api

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"golang.org/x/crypto/bcrypt"
)

const (
	INSERT_SESSION         = "INSERT INTO session (id, user_id, expires_at) VALUES (:id, :user_id, NOW() + interval '60 minutes');"
	GET_USER_BY_SESSION_ID = "SELECT user_info.* FROM session JOIN user_info ON session.user_id = user_info.id WHERE session.id=$1 AND expires_at >= now();"
)

type Session struct {
	Id        string    `json:"id" db:"id"`
	UserId    int       `json:"userId" db:"user_id"`
	ExpiresAt time.Time `json:"expiresAt" db:"expires_at"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func NewSessionFromUser(u *User) *Session {
	return &Session{
		Id:     uuid.New().String(),
		UserId: u.Id,
	}
}

func LoadUserFromSessionId(id string, db *sqlx.DB) (*User, error) {
	u := &User{}
	err := db.Get(u, GET_USER_BY_SESSION_ID, id)
	if err != nil {
		return nil, err
	}
	return u, nil
}

func (s *Session) CreateInDB(db *sqlx.DB) error {
	_, err := db.NamedExec(INSERT_SESSION, s)
	return err
}

func Auth(c *gin.Context) {
	l := &LoginRequest{}
	//vérifier que les identifiants sont présents
	if err := c.BindJSON(l); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid request, please provide valid email and password",
		})
		return
	}
	u, err := LoadUserFromEmail(APIDatabase, l.Email)
	//vérifier que l'email est bien attribuée à un utilisateur
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "cannot find user with provided credentials, please try again",
		})
		return
	}
	//Vérifier que le mot de passe est correct
	err = bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(l.Password))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "incorrect email or password, please try again",
		})
		return
	}
	s := NewSessionFromUser(u)
	//vérifier que la session a bien été stockée dans la bdd
	if err = s.CreateInDB(APIDatabase); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "error has occured while creating session, please try again later",
		})
		return
	}
	c.JSON(http.StatusCreated, gin.H{
		"sessID": s.Id,
		"userID": u.Id,
	})
}
