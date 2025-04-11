package categories

import (
	"database/sql"
	"net/http"
	dbase "platifyapi/db"
	"platifyapi/models"

	"github.com/gin-gonic/gin"
)

func insertCategory(db *sql.DB, category *models.Category) (*models.Category, error) {
	query := `INSERT INTO categories (name, description) VALUES ($1, $2) RETURNING category_id`
	err := db.QueryRow(query, category.Name, category.Description).Scan(&category.CategoryID)
	if err != nil {
		return nil, err
	}

	return category, nil
}

func InsertCategory(c *gin.Context) {
	db, err := dbase.ConnectDB()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer db.Close()

	var category models.Category
	if err := c.ShouldBindJSON(&category); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}	
	
	insertedCategory, err := insertCategory(db, &category)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, insertedCategory)
}

func getCategories(db *sql.DB) ([]models.Category, error) {
	query := `SELECT category_id, name, description FROM categories`
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categories []models.Category
	for rows.Next() {
		var category models.Category
		if err := rows.Scan(&category.CategoryID, &category.Name, &category.Description); err != nil {
			return nil, err
		}
		categories = append(categories, category)
	}

	return categories, nil
}

func GetCategories(c *gin.Context){
	db, err := dbase.ConnectDB()
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    defer db.Close()

	categories, err := getCategories(db)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, categories)
}

