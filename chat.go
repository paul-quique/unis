package api

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

const (
	SELECT_CONVERSATION  = "SELECT * FROM message WHERE (sent_from=$1 AND sent_to=$2) OR (sent_from=$2 AND sent_to=$1) ORDER BY sent_at"
	SELECT_CONVERSATIONS = "SELECT DISTINCT sent_to FROM message WHERE sent_from = $1 UNION DISTINCT SELECT DISTINCT sent_from FROM message WHERE sent_to = $1;"
	INSERT_MESSAGE       = "INSERT INTO message (sent_from, sent_to, sent_at, content) VALUES (:sent_from, :sent_to, NOW(), :content);"
)

type Message struct {
	Id       int       `json:"id" db:"id"`
	SentFrom int       `json:"sentFrom" db:"sent_from"`
	SentTo   int       `json:"sentTo" db:"sent_to"`
	Content  string    `json:"content" db:"content"`
	SentAt   time.Time `json:"sentAt" db:"sent_at"`
}

type PostMessageRequest struct {
	SessionId string `json:"sessId"`
	SentTo    int    `json:"sentTo" db:"send_to"`
	Content   string `json:"content" db:"content"`
}

type PostConversationRequest struct {
	SessionId string `json:"sessId"`
	SentTo    int    `json:"sentTo" db:"send_to"`
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
	c.Status(http.StatusCreated)
}

func PostConversation(c *gin.Context) {
	//extraire les paramètres dans un struct pour vérifier qu'ils sont valides
	cr := &PostConversationRequest{}
	err := c.BindJSON(cr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "cannot bind json, please check the request is valid",
		})
		return
	}
	//vérifier que la session n'est pas expirée
	u, err := LoadUserFromSessionId(cr.SessionId, APIDatabase)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "session expired, please re-login and try again",
		})
		return
	}
	//vérifier que l'émetteur est différent du destinataire
	if u.Id == cr.SentTo {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "sender and receiver should be different",
		})
		return
	}
	//obtenir les messages demandés
	messages := &[]*Message{}
	err = APIDatabase.Select(messages, SELECT_CONVERSATION, u.Id, cr.SentTo)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	//renvoyer les messages de la conversation demandée
	c.JSON(http.StatusOK, messages)
}

func PostConversations(c *gin.Context) {
	//extraire les paramètres dans un struct pour vérifier qu'ils sont valides
	ar := &AuthenticatedRequest{}
	err := c.BindJSON(ar)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "cannot bind json, please check the request is valid",
		})
		return
	}
	//vérifier que la session n'est pas expirée
	u, err := LoadUserFromSessionId(ar.SessionId, APIDatabase)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "session expired, please re-login and try again",
		})
		return
	}
	ids := &[]int{}
	//renvoyer les ids des utilisateurs qui ont une conversation avec
	//l'auteur de la requête
	err = APIDatabase.Select(ids, SELECT_CONVERSATIONS, u.Id)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, ids)
}
