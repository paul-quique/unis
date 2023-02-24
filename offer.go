package api

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

const (
	ADD_POINTS_TO_USER    = "UPDATE user_info SET points = points + $1 WHERE id = $2"
	REMOVE_POINTS_TO_USER = "UPDATE user_info SET points = points - $1 WHERE id = $2"
	INSERT_TRANSACTION    = "INSERT INTO transaction (offer_id, return_date) VALUES (:offer_id, NOW() + interval '7 day');"
	INSERT_OFFER          = "INSERT INTO offer (borrower_id, lender_id, product_id, created_at, expires_at) VALUES (:borrower_id, :lender_id, :product_id, NOW(), NOW() + interval '7 day');"
	GET_OFFER_BY_ID       = "SELECT * FROM offer WHERE id=$1;"
	GET_OFFERS_BY_ID      = "SELECT * FROM offer WHERE borrower_id=$1 OR lender_id=$1 AND id NOT IN (SELECT offer_id FROM transaction);"
)

type Offer struct {
	Id         int       `json:"id" db:"id"`
	BorrowerId int       `json:"borrowerId" db:"borrower_id"`
	LenderId   int       `json:"lenderId" db:"lender_id"`
	ProductId  int       `json:"productId" db:"product_id"`
	CreatedAt  time.Time `json:"createdAt" db:"created_at"`
	ExpiresAt  time.Time `json:"expiresAt" db:"expires_at"`
}

type Transaction struct {
	Id         int       `json:"id" db:"id"`
	OfferId    int       `json:"offerId" db:"offer_id"`
	ReturnDate time.Time `json:"returnDate" db:"return_date"`
}

type CreateOfferRequest struct {
	SessionId string `json:"sessId"`
	ProductId int    `json:"productId" db:"product_id"`
}

type AcceptOfferRequest struct {
	OfferId   int    `json:"offerId" db:"offer_id"`
	SessionID string `json:"sessId"`
}

type GetOffersRequest struct {
	SessionID string `json:"sessId"`
}

func (o *Offer) CreateInDB(db *sqlx.DB) (err error) {
	_, err = db.NamedExec(INSERT_OFFER, o)
	return err
}

func (t *Transaction) CreateInDB(db *sqlx.DB) (err error) {
	_, err = db.NamedExec(INSERT_TRANSACTION, t)
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
	u, err := LoadUserFromSessionId(or.SessionId, APIDatabase)
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

func LoadOfferFromId(id int, db *sqlx.DB) (*Offer, error) {
	o := &Offer{}
	err := db.Get(o, GET_OFFER_BY_ID, id)
	if err != nil {
		return nil, err
	}
	return o, nil
}

func AcceptOffer(c *gin.Context) {
	//extraire les paramètres dans un struct pour vérifier qu'ils sont valides
	or := &AcceptOfferRequest{}
	err := c.BindJSON(or)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "cannot bind json, please check the request is valid",
		})
		return
	}
	//obtenir l'offre que l'utilisateur souhaite accepter
	o, err := LoadOfferFromId(or.OfferId, APIDatabase)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "please specify a valid offer id",
		})
		return
	}
	//vérifiez que l'auteur de la requête est bien celui qui a reçu l'offre
	author, err := LoadUserFromSessionId(or.SessionID, APIDatabase)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid session, re-login and try again",
		})
		return
	}
	if author.Id != o.LenderId {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "you must own this offer to accept it",
		})
		return
	}
	//obtenir le produit sur lequel porte l'offre
	p, err := LoadProductFromId(APIDatabase, o.ProductId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "error, cannot find the product",
		})
		return
	}
	//vérifier que l'emprunteur a les fonds nécessaires pour louer le produits
	borr, err := LoadUserFromId(APIDatabase, o.BorrowerId)
	if borr.Points < p.Price {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "error, the borrower doesn't have enough points",
		})
		return
	}
	//actualiser les points des comptes qui participent à la transaction
	_, err = APIDatabase.Exec(REMOVE_POINTS_TO_USER, p.Price, o.BorrowerId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "error, please contact an administrator",
		})
		return
	}
	_, err = APIDatabase.Exec(ADD_POINTS_TO_USER, p.Price, o.BorrowerId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "error, please contact an administrator",
		})
		return
	}
	//créer une transaction pour le produit
	t := &Transaction{
		OfferId: o.Id,
	}
	//enregistrer la transaction créée dans la base de données
	err = t.CreateInDB(APIDatabase)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "cannot save transaction in database, please contact an administrator",
		})
		return
	}
	c.Status(http.StatusOK)
}

func GetOffers(c *gin.Context) {
	//extraire les paramètres dans un struct pour vérifier qu'ils sont valides
	gor := &GetOffersRequest{}
	err := c.BindJSON(gor)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "cannot bind json, please check the request is valid",
		})
		return
	}

	//vérifier que la session n'est pas expirée
	u, err := LoadUserFromSessionId(gor.SessionID, APIDatabase)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "session expired, please re-login and try again",
		})
		return
	}

	//faire une requête pour obtenir l'ensemble des offres effectuées ou
	//bien reçues par l'utilisateur
	o := &[]Offer{}
	err = APIDatabase.Select(o, GET_OFFERS_BY_ID, u.Id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "error while loading offers from database, please contact an administrator",
		})
		return
	}

	//return offers to the client
	c.JSON(http.StatusOK, o)
}
