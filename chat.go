package api

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

const (
	INSERT_MESSAGE = "INSERT INTO message (sent_from, sent_to, sent_at, content) VALUES (:sent_from, :sent_to, NOW(), :content);"
)

type Message struct {
	Id       int       `json:"id" db:"id"`
	SentFrom int       `json:"sentFrom" db:"sent_from"`
	SentTo   int       `json:"sentTo" db:"send_to"`
	Content  string    `json:"content" db:"content"`
	SentAt   time.Time `json:"sentAt" db:"sent_at"`
}

type PostMessageRequest struct {
	SessionId string `json:"sessId"`
	SentTo    int    `json:"sentTo" db:"send_to"`
	Content   string `json:"content" db:"content"`
}

func (m *Message) CreateInDB(db *sqlx.DB) error {
	_, err := db.NamedExec(INSERT_MESSAGE, m)
	return err
}

func PostMessage(c *gin.Context) {
	//extraire les paramètres dans un struct pour vérifier qu'ils sont valides
	mr := &PostMessageRequest{}
	err := c.BindJSON(mr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "cannot bind json, please check the request is valid",
		})
		return
	}
	//vérifier que la session n'est pas expirée
	u, err := LoadUserFromSessionId(mr.SessionId, APIDatabase)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "session expired, please re-login and try again",
		})
		return
	}
	//vérifier que l'émetteur est différent du destinataire
	if u.Id == mr.SentTo {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "sender and receiver should be different",
		})
		return
	}
	m := &Message{SentFrom: u.Id, SentTo: mr.SentTo, Content: mr.Content}
	//enregistrer le message dans la base de données
	err = m.CreateInDB(APIDatabase)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
}
