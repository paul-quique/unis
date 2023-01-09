package api

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

const (
	INSERT_OFFER = "INSERT INTO offer (borrower_id, lender_id, product_id, created_at, expires_at) VALUES (:borrower_id, :lender_id, :product_id, NOW(), NOW() + interval '7 day');"
)

type Offer struct {
	Id         int       `json:"id" db:"id"`
	BorrowerId int       `json:"borrowerId" db:"borrower_id"`
	LenderId   int       `json:"lenderId" db:"lender_id"`
	ProductId  int       `json:"productId" db:"product_id"`
	CreatedAt  time.Time `json:"createdAt" db:"created_at"`
	ExpiresAt  time.Time `json:"expiresAt" db:"expires_at"`
}

type CreateOfferRequest struct {
	SessionID string `json:"sessID"`
	ProductId int    `json:"productId" db:"product_id"`
}

func (o *Offer) CreateInDB(db *sqlx.DB) (err error) {
	_, err = db.NamedExec(INSERT_OFFER, o)
	return err
}

func PostOffer(c *gin.Context) {
	//extraire les paramètres dans un struct pour vérifier qu'ils sont valides
	or := &CreateOfferRequest{}
	err := c.BindJSON(or)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "cannot bind json, please check the request is valid",
		})
		return
	}
	//charger le produit depuis la base de données
	p, err := LoadProductFromId(APIDatabase, or.ProductId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "cannot make an offer on this product, please check product id",
		})
		return
	}
	//vérifier que la session n'est pas expirée
	u, err := LoadUserFromSessionId(or.SessionID, APIDatabase)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "session expired, please re-login and try again",
		})
		return
	}
	//vérifier que l'emprunteur n'est pas le prêteur
	if p.UserId == u.Id {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "lender and borrower are identical",
		})
		return
	}
	//vérifier que la création s'est bien effectuée
	o := &Offer{
		BorrowerId: u.Id,
		LenderId:   p.UserId,
		ProductId:  p.Id,
	}
	err = o.CreateInDB(APIDatabase)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "error creating offer, please try again later",
		})
		return
	}
	//l'offre a été créé avec succès
	c.Status(http.StatusCreated)
}
