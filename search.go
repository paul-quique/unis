package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

const (
	GET_SEARCH_RESULTS = "SELECT * FROM product ORDER BY similarity(product.name, $1) DESC LIMIT $2 OFFSET $3;"
)

type SearchProductRequest struct {
	Query        string `json:"query"`
	Offset       int    `json:"offset"`
	ItemsPerPage int    `json:"itemsPerPage"`
}

func LoadSearchResults(db *sqlx.DB, s *SearchProductRequest) (*[]Product, error) {
	p := &[]Product{}
	err := db.Select(p, GET_SEARCH_RESULTS, s.Query, s.ItemsPerPage, s.Offset)
	if err != nil {
		return nil, err
	}
	return p, nil
}

func PostSearch(c *gin.Context) {
	//extraire les paramètres dans un struct pour vérifier qu'ils sont valides
	sr := &SearchProductRequest{}
	err := c.BindJSON(sr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "cannot bind json, please check the request is valid",
		})
		return
	}
	//récupérer les produits dans la base de données par pertinence
	//en prenant en compte l'offset et le nombre de produits demandés
	p, err := LoadSearchResults(APIDatabase, sr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	//renvoyer les produits correspondants
	c.JSON(http.StatusOK, p)
}
