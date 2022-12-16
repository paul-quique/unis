package api

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

const (
	MaxProductNameLength = 45
	MinProductNameLength = 2
	INSERT_PRODUCT       = "INSERT INTO product (name, category_id, price, user_id) VALUES (:name, :category_id, :price, :user_id);"
)

type Product struct {
	Id         int    `json:"id" db:"id"`
	Name       string `json:"name" db:"name"`
	CategoryId int    `json:"categoryId" db:"category_id"`
	Price      int    `json:"price" db:"price"`
	UserId     int    `json:"userId" db:"user_id"`
}

type CreateProductRequest struct {
	SessionID  string `json:"sessID"`
	Name       string `json:"name" db:"name"`
	CategoryId int    `json:"categoryId" db:"category_id"`
	Price      int    `json:"price" db:"price"`
}

func (p *CreateProductRequest) IsValid() string {
	//vérifier que le nom a entre 4 et 32 caractères
	l := len(p.Name)
	if !(l >= 4 && l <= 32) {
		return "name length must be between 4 and 32 chars"
	}
	//vérifier que la catégorie spécifiée existe
	_, err := LoadCategoryFromId(APIDatabase, p.CategoryId)
	if err != nil {
		return "provided category id doesn't exist"
	}
	//vérifier que le prix n'est pas nul ou négatif
	if p.Price <= 0 {
		return "price must be positive"
	}
	return "ok"
}
func (p *Product) CreateInDB(db *sqlx.DB) error {
	_, err := db.NamedExec(INSERT_PRODUCT, p)
	return err
}

func PostProduct(c *gin.Context) {
	//extraire les paramètres dans un struct pour vérifier qu'ils sont valides
	pr := &CreateProductRequest{}
	err := c.BindJSON(pr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "cannot bind json, please check the request is valid",
		})
		return
	}
	if m := pr.IsValid(); m != "ok" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": m,
		})
		return
	}
	//vérifier que la session n'est pas expirée
	u, err := LoadUserFromSessionId(pr.SessionID, APIDatabase)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "session expired, please re-login and try again",
		})
		return
	}
	//créer un produit et l'enregistrer dans la base de données
	p := &Product{
		Name:       pr.Name,
		UserId:     u.Id,
		Price:      pr.Price,
		CategoryId: pr.CategoryId,
	}
	//vérifier que la création s'est bien effectuée
	err = p.CreateInDB(APIDatabase)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "error creating product, please try again later",
		})
		return
	}
	//le produit a été créé avec succès
	c.JSON(http.StatusCreated, gin.H{
		"createdAt": time.Now(),
	})
}
