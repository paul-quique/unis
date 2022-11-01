package api

import (
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

const (
	INSERT_SESSION         = "INSERT INTO session (id, user_id, expires_at) VALUES (:id, :user_id, :expires_at);"
	GET_USER_BY_SESSION_ID = "SELECT user_info.* FROM session JOIN user_info ON session.user_id = user_info.id WHERE session.id=$1 AND expires_at <= now();"
)

type Session struct {
	Id        string    `json:"id" db:"id"`
	UserId    int       `json:"userId" db:"user_id"`
	ExpiresAt time.Time `json:"expiresAt" db:"expires_at"`
}

func NewSessionFromUser(u *User) *Session {
	return &Session{
		Id:        uuid.New().String(),
		UserId:    u.Id,
		ExpiresAt: time.Now().Add(60 * time.Minute),
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
