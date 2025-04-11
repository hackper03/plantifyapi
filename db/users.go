package dbase

import (
	"database/sql"
	"fmt"
	"net/http"
	"platifyapi/models"
	"platifyapi/util"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
)

func insertUser(db *sql.DB, user *models.User) (*models.User, error) {
	fmt.Printf("Internal User: %v\n", user)
	query := `
        INSERT INTO users (name, email, password, phone, address, created_at, updated_at) 
        VALUES ($1, $2, $3, $4, $5, $6, $7)
        RETURNING user_id;
    `
	stmt, err := db.Prepare(query)
	if err != nil {
		return nil, fmt.Errorf("could not prepare statement: %v", err)
	}
	defer stmt.Close()

	fmt.Printf("User before hashing: %v\n", user)
	hashedPassword, err := util.HashPassword(user.Password)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %v", err)
	}

	fmt.Printf("Hash generated %v", hashedPassword)

	err = stmt.QueryRow(user.Name, user.Email, hashedPassword, user.Phone, user.Address, time.Now(), time.Now()).Scan(&user.UserID)
	if err != nil {
		// Check if the error is a unique constraint violation
		if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == "23505" {
			return nil, fmt.Errorf("email '%s' already exists", user.Email)
		}
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("email '%s' already exists", user.Email)
		}
		return nil, fmt.Errorf("could not store user: %v", err)
	}

	return user, nil
}

// Handle the Create User Request
func CreateUser(c *gin.Context) {
	db, err := ConnectDB()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "database connection error", "details": err.Error()})
		return
	}
	defer db.Close()

	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body", "details": err.Error()})
		return
	}

	savedUser, err := insertUser(db, &user)
	if err != nil {
		if strings.Contains(err.Error(), "email already exists") {
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create user", "details": err.Error()})
		return
	}

	token, err := util.GenerateToken(savedUser.Email, savedUser.UserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "could not generate token",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "User created successfully",
		"token":   token,
	})
}

func Login(c *gin.Context) {
	db, err := ConnectDB()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Error while connecting Database": err.Error()})
		return
	}
	defer db.Close()

	var user models.Login
	if err = c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Bindding Error",
			"error":   err.Error(),
		})
	}
	err = validateUser(db, &user)
	if err != nil {
		if err.Error() == "email not found" {
			c.JSON(http.StatusNotFound, gin.H{
				"message": "user not found",
			})
			return
		} else if err.Error() == "invalid password" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"message": "Invalid Password Credentials",
			})
			return
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "Internal Server Error",
			})
			return
		}
	}

	token, err := util.GenerateToken(user.Email, user.UserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Could not fetch token from jwt",
			"error":   err.Error(),
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Login Successfully",
		"token":   token,
	})
}

func validateUser(db *sql.DB, user *models.Login) error {
	query := `select user_id, email, password from users where email = $1;`
	row := db.QueryRow(query, user.Email)
	var storedEmail, storedPassword string
	err := row.Scan(&user.UserID, &storedEmail, &storedPassword)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("email not found")
		}
		return fmt.Errorf("error querying user: %v", err)
	}

	if !util.CheckPasswordHash(user.Password, storedPassword) {
		return fmt.Errorf("invalid password")
	}

	return nil
}
