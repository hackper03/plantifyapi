package cart

import (
	"database/sql"
	"fmt"
	"net/http"
	dbase "platifyapi/db"
	"platifyapi/models"
	"platifyapi/util"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// Insert a new cart into the database
func insertCart(db *sql.DB, userID int) (int, error) {
	query := `INSERT INTO cart (user_id, created_at, updated_at) VALUES ($1, $2, $3) RETURNING cart_id`
	var cartID int
	err := db.QueryRow(query, userID, time.Now(), time.Now()).Scan(&cartID)
	if err != nil {
		return 0, err
	}
	return cartID, nil
}

// Create a new cart
func CreateCart(c *gin.Context) {
	db, err := dbase.ConnectDB()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error connecting to database"})
		return
	}
	defer db.Close()

	var cart models.Cart
	if err := c.ShouldBindJSON(&cart); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token, err := c.Cookie("authToken")
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	fmt.Printf("Token: %s\n", token)

	userID, err := util.VerifyToken(token)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		return
	}
	cart.UserID = *userID
	cartID, err := insertCart(db, int(cart.UserID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Cart created successfully",
		"cart_id": cartID,
	})
}

// Get all carts
func getCarts(db *sql.DB) ([]models.Cart, error) {
	query := `SELECT cart_id, user_id, created_at, updated_at FROM cart`
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var carts []models.Cart
	for rows.Next() {
		var cart models.Cart
		if err := rows.Scan(&cart.CartID, &cart.UserID, &cart.CreatedAt, &cart.UpdatedAt); err != nil {
			return nil, err
		}
		carts = append(carts, cart)
	}

	return carts, nil
}

func getCartByID(db *sql.DB, cartID int) (*models.Cart, error) {
	query := `SELECT cart_id, user_id, created_at, updated_at FROM cart WHERE cart_id = $1`
	row := db.QueryRow(query, cartID)

	var cart models.Cart
	if err := row.Scan(&cart.CartID, &cart.UserID, &cart.CreatedAt, &cart.UpdatedAt); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // No cart found with the given ID
		}
		return nil, err
	}

	return &cart, nil
}

// Get cart by ID handler
func GetCartByID(c *gin.Context) {
	db, err := dbase.ConnectDB()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error connecting to database"})
		return
	}
	defer db.Close()

	cartID := c.Param("id")
	id, err := strconv.Atoi(cartID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid cart ID"})
		return
	}

	cart, err := getCartByID(db, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if cart == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Cart not found"})
		return
	}
	c.JSON(http.StatusOK, cart)
}

// Get all carts handler
func GetCarts(c *gin.Context) {
	db, err := dbase.ConnectDB()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error connecting to database"})
		return
	}
	defer db.Close()

	carts, err := getCarts(db)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, carts)
}

// Insert a new cart item into the database
func insertCartItem(db *sql.DB, cartItem *models.CartItem) (int, error) {
	query := `INSERT INTO cart_items (cart_id, plant_id, service_id, quantity, price) 
    VALUES ($1, $2, $3, $4, $5) RETURNING cart_item_id`
	var cartItemID int
	err := db.QueryRow(query, cartItem.CartID, cartItem.PlantID, cartItem.ServiceID, cartItem.Quantity, cartItem.Price).Scan(&cartItemID)
	if err != nil {
		return 0, err
	}
	return cartItemID, nil
}

// Create a new cart item
func CreateCartItem(c *gin.Context) {
	db, err := dbase.ConnectDB()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error connecting to database"})
		return
	}
	defer db.Close()

	var cartItem models.CartItem
	if err := c.ShouldBindJSON(&cartItem); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	cartItemID, err := insertCartItem(db, &cartItem)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message":      "Cart item created successfully",
		"cart_item_id": cartItemID,
	})
}

// Get all cart items for a specific cart
func getCartItemsByCartID(db *sql.DB, cartID int) ([]models.CartItem, error) {
	query := `SELECT cart_item_id, cart_id, plant_id, service_id, quantity, price, total 
    FROM cart_items WHERE cart_id = $1`
	rows, err := db.Query(query, cartID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var cartItems []models.CartItem
	for rows.Next() {
		var cartItem models.CartItem
		if err := rows.Scan(&cartItem.CartItemID, &cartItem.CartID, &cartItem.PlantID, &cartItem.ServiceID, &cartItem.Quantity, &cartItem.Price, &cartItem.Total); err != nil {
			return nil, err
		}
		cartItems = append(cartItems, cartItem)
	}

	return cartItems, nil
}

// Get cart items by cart ID handler
func GetCartItemsByCartID(c *gin.Context) {
	db, err := dbase.ConnectDB()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error connecting to database"})
		return
	}
	defer db.Close()

	cartID := c.Param("id")
	id, err := strconv.Atoi(cartID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid cart ID"})
		return
	}

	cartItems, err := getCartItemsByCartID(db, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, cartItems)
}
