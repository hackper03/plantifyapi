package users

import (
	"database/sql"
	"fmt"
	"net/http"
	dbase "platifyapi/db"
	"platifyapi/models"
	"platifyapi/util"
	"time"

	"github.com/gin-gonic/gin"
)

func insertUser(db *sql.DB, user *models.User) (int64, error) {
	query := `
        INSERT INTO users (name, email, password, phone, address, created_at, updated_at) 
        VALUES ($1, $2, $3, $4, $5, $6, $7)
        RETURNING user_id;
    `
	stmt, err := db.Prepare(query)
	if err != nil {
		fmt.Printf("Could not prepare statement")
		panic(err)
	}
	defer stmt.Close()

	hashedPassword, err := util.HashPassword(user.Password)
	if err != nil {
		fmt.Printf("failee to convert password into hash")
		panic(err)
	}

	err = stmt.QueryRow(user.Name, user.Email, hashedPassword, user.Phone, user.Address, time.Now(), time.Now()).Scan(&user.UserID)
	if err != nil {
		fmt.Printf("Could not store user: %v", err)
		panic(err)
	}

	return user.UserID, nil
}

// Handle the Create User Request
func CreateUser(c *gin.Context) {
	db, err := dbase.ConnectDB()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Database connection error",
			"error":   err.Error(),
			"success": false,
		})
		return
	}
	defer db.Close()

	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Binding Error",
			"error":   err.Error(),
			"success": false,
		})
		return
	}

	_, err = insertUser(db, &user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Could not store user",
			"error":   err.Error(),
			"success": false,
		})
		return
	}

	token, err := util.GenerateToken(user.Email, user.UserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Could not fetch token from jwt",
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	c.SetCookie(
		"authToken", // name
		token,       // value
		3600,        // max age in seconds (1 hour)
		"/",         // path
		"",          // domain
		false,       // secure
		true,        // httpOnly
	)

	c.JSON(http.StatusCreated, gin.H{
		"message": "User created succesfully",
		"success": true,
	})
}

func Login(c *gin.Context) {
	db, err := dbase.ConnectDB()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Error while connecting Database": err.Error()})
		return
	}
	defer db.Close()

	var user models.Login
	if err = c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Binding Error",
			"error":   err.Error(),
			"success": false,
		})
		return
	}

	err = validateUser(db, &user)
	if err != nil {
		if err.Error() == "email not found" {
			c.JSON(http.StatusNotFound, gin.H{
				"message": "user not found",
				"success": false,
			})
			return
		} else if err.Error() == "invalid password" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"message": "Invalid Password Credentials",
				"success": false,
			})
			return
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "Internal Server Error",
				"success": false,
			})
			return
		}
	}

	token, err := util.GenerateToken(user.Email, user.UserID)
	if err != nil {
		fmt.Printf("Failed to fetch token")
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Could not fetch token from jwt",
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	// Set the token in an HTTP-only cookie
	c.SetCookie(
		"authToken", // name
		token,       // value
		3600,        // max age in seconds (1 hour)
		"/",         // path
		"",          // domain
		false,       // secure
		true,        // httpOnly
	)

	c.JSON(http.StatusOK, gin.H{
		"message": "Login Successfully",
		"success": true,
	})
}

func Logout(c *gin.Context) {
	// Clear the auth cookie by setting it to empty and expiring it immediately
	c.SetCookie(
		"authToken", // name
		"",          // empty value
		-1,          // negative max age to expire immediately
		"/",         // path
		"",          // domain
		false,       // secure
		true,        // httpOnly
	)

	c.JSON(http.StatusOK, gin.H{
		"message": "Logged out successfully",
		"success": true,
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
