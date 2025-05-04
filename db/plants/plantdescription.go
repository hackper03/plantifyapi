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

func insertPlantDescription(db *sql.DB, description *models.PlantDescription) (*models.PlantDescription, error) {
	query := `INSERT INTO plant_description (plant_id, description, features, care_instructions, created_at, updated_at) 
    VALUES ($1, $2, $3, $4, $5, $6) RETURNING plant_id`
	err := db.QueryRow(query, description.PlantID, description.Description, description.Features, description.CareInstructions, time.Now(), time.Now()).Scan(&description.PlantID)
	if err != nil {
		return nil, err
	}

	return description, nil
}

func InsertPlantDescription(c *gin.Context) {
	db, err := dbase.ConnectDB()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Error while connecting Database": err.Error()})
		return
	}
	defer db.Close()

	var description models.PlantDescription
	if err := c.ShouldBindJSON(&description); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Insert the plant description into the database
	_, err = insertPlantDescription(db, &description)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, description)
}

func getPlantDescriptionByID(db *sql.DB, id int) (*models.PlantDescription, error) {
	query := `SELECT plant_id, description, features, care_instructions FROM plant_description WHERE plant_id = $1`
	row := db.QueryRow(query, id)

	var description models.PlantDescription
	if err := row.Scan(&description.PlantID, &description.Description, &description.Features, &description.CareInstructions); err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("no description found for plant ID %d", id)
		}
		return nil, err
	}

	return &description, nil
}

func GetPlantDescriptionByID(c *gin.Context) {
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

	description, err := getPlantDescriptionByID(db, id)
	if err != nil {
		if err.Error() == fmt.Sprintf("no description found for plant ID %d", id) {
			c.JSON(http.StatusNotFound, gin.H{"error": "No description found for the given plant ID"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, description)
}
