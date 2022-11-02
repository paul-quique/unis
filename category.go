package api

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

const (
	GET_CATEGORIES     = "SELECT * FROM category;"
	GET_CATEGORY_BY_ID = "SELECT * FROM category WHERE id=$1;"
	INSERT_CATEGORY    = "INSERT INTO category (name) VALUES (:name);"
)

type Category struct {
	Id   int    `json:"id" db:"id"`
	Name string `json:"name" db:"name"`
}

func (c *Category) CreateInDb(db *sqlx.DB) error {
	_, err := db.NamedExec(INSERT_CATEGORY, c)
	return err
}

func LoadCategoryFromId(db *sqlx.DB, id int) (*Category, error) {
	c := &Category{}
	err := db.Get(c, GET_CATEGORY_BY_ID, id)
	if err != nil {
		return nil, err
	}
	return c, nil
}

func LoadCategories(db *sqlx.DB) ([]*Category, error) {
	c := []*Category{}
	err := db.Get(c, GET_CATEGORIES)
	if err != nil {
		return nil, err
	}
	return c, nil
}

func GetCategory(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "please provide a valid category id",
		})
		return
	}

	cat, err := LoadCategoryFromId(APIDatabase, id)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "cannot find category with the provided id",
		})
		return
	}

	c.JSON(http.StatusOK, cat)
}

func GetCategories(c *gin.Context) {
	cats, err := LoadCategories(APIDatabase)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "cannot load categories, please try again later",
		})
		return
	}
	c.JSON(http.StatusOK, cats)
}
