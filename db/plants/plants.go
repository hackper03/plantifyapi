package plants

import (
	"database/sql"
	"fmt"
	"net/http"
	dbase "platifyapi/db"
	"platifyapi/models"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

func insertPlant(db *sql.DB, plant *models.Plant) (*models.Plant, error) {
	query := `INSERT INTO plants (name, description, category_id, image_url, rating, price, original_price, badge_text, created_at, updated_at) 
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10) RETURNING plant_id`
	err := db.QueryRow(query, plant.Name, plant.Description, plant.CategoryID, plant.ImageURL, plant.Rating, plant.Price, plant.OriginalPrice, plant.BadgeText, time.Now(), time.Now()).Scan(&plant.PlantID)
	if err != nil {
		return nil, err
	}

	return plant, nil
}

func InsertPlant(c *gin.Context) {
	db, err := dbase.ConnectDB()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Error while connecting Database": err.Error()})
		return
	}
	defer db.Close()

	var plant models.Plant
	if err := c.ShouldBindJSON(&plant); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Insert the plant into the database
	_, err = insertPlant(db, &plant)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, plant)
}

func getPlants(db *sql.DB) ([]models.Plant, error) {
	query := `SELECT plant_id, name, description, category_id, image_url, rating, price, original_price, badge_text FROM plants`
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var plants []models.Plant
	for rows.Next() {
		var plant models.Plant
		if err := rows.Scan(&plant.PlantID, &plant.Name, &plant.Description, &plant.CategoryID, &plant.ImageURL, &plant.Rating, &plant.Price, &plant.OriginalPrice, &plant.BadgeText); err != nil {
			return nil, err
		}
		plants = append(plants, plant)
	}

	return plants, nil
}

func getPlantsByCateogryID(db *sql.DB, categoryID int) ([]models.Plant, error) {
	query := `SELECT plant_id, name, description, category_id, image_url, rating, price, original_price, badge_text FROM plants WHERE category_id = $1`
	rows, err := db.Query(query, categoryID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var plants []models.Plant
	for rows.Next() {
		var plant models.Plant
		if err := rows.Scan(&plant.PlantID, &plant.Name, &plant.Description, &plant.CategoryID, &plant.ImageURL, &plant.Rating, &plant.Price, &plant.OriginalPrice, &plant.BadgeText); err != nil {
			return nil, err
		}
		plants = append(plants, plant)
	}

	if len(plants) == 0 {
		return nil, fmt.Errorf("no plants found for category ID %d", categoryID)
	}

	return plants, nil
}

func GetPlantsByCateogryID(c *gin.Context) {
	db, err := dbase.ConnectDB()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer db.Close()

	categoryID := c.Param("id")
	id, err := strconv.Atoi(categoryID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid plant ID"})
		return
	}
	plants, err := getPlantsByCateogryID(db, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if len(plants) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "No plants found for the given category ID"})
		return
	}

	c.JSON(http.StatusOK, plants)
}
func getPlantByID(db *sql.DB, id int) (*models.Plant, error) {
	query := `SELECT plant_id, name, description, category_id, image_url, rating, price, original_price, badge_text FROM plants WHERE plant_id = $1`
	row := db.QueryRow(query, id)

	var plant models.Plant
	if err := row.Scan(&plant.PlantID, &plant.Name, &plant.Description, &plant.CategoryID, &plant.ImageURL, &plant.Rating, &plant.Price, &plant.OriginalPrice, &plant.BadgeText); err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("no plant found with ID %d", id)
		}
		return nil, err
	}

	return &plant, nil
}

func GetPlantByID(c *gin.Context) {
	db, err := dbase.ConnectDB()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer db.Close()

	plantID := c.Param("id")
	id, err := strconv.Atoi(plantID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid plant ID"})
		return
	}

	plant, err := getPlantByID(db, id)
	if err != nil {
		if err.Error() == "no plant found with ID" {
			c.JSON(http.StatusNotFound, gin.H{"error": "No plant found with the given ID"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, plant)
}

func GetPlants(c *gin.Context) {
	db, err := dbase.ConnectDB()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer db.Close()

	plants, err := getPlants(db)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, plants)
}
