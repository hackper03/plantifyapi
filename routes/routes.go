package routes

import (
	"platifyapi/db/cart"
	"platifyapi/db/categories"
	"platifyapi/db/plants"
	"platifyapi/db/users"
	"strings"

	"slices"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(server *gin.Engine) {
	server.Use(corsMiddleware())
	server.GET("/category", categories.GetCategories)
	server.POST("/category", categories.InsertCategory)
	server.GET("/plants", plants.GetPlants)
	server.POST("/plant", plants.InsertPlant)
	server.GET("/plant/category/:id", plants.GetPlantsByCateogryID)
	server.POST("/signup", users.CreateUser)
	server.POST("/login", users.Login)
	server.POST("/logout", users.Logout)
	server.GET("/plant/:id", plants.GetPlantByID)
	server.POST("/api/cart", cart.CreateCart)
	server.GET("api/cart-items/:id", cart.GetCartItemsByCartID)
	server.POST("/api/cart-items", cart.CreateCartItem)
	server.GET("/api/cart/:id", cart.GetCartByID)
	server.GET("/api/plantdescription/:id", plants.GetPlantDescriptionByID)
	server.POST("/api/plantdescription", plants.InsertPlantDescription)
}

// // CORS middleware function definition
// func corsMiddleware() gin.HandlerFunc {
// 	// Define allowed origins as a comma-separated string
// 	originsString := "http://localhost:5173/"
// 	var allowedOrigins []string
// 	if originsString != "" {
// 		// Split the originsString into individual origins and store them in allowedOrigins slice
// 		allowedOrigins = strings.Split(originsString, ",")
// 	}

// 	// Return the actual middleware handler function
// 	return func(c *gin.Context) {
// 		// Function to check if a given origin is allowed
// 		isOriginAllowed := func(origin string, allowedOrigins []string) bool {
// 			for _, allowedOrigin := range allowedOrigins {
// 				if origin == allowedOrigin {
// 					return true
// 				}
// 			}
// 			return false
// 		}

// 		// Get the Origin header from the request
// 		origin := c.Request.Header.Get("Origin")

// 		// Check if the origin is allowed
// 		if isOriginAllowed(origin, allowedOrigins) {
// 			// If the origin is allowed, set CORS headers in the response
// 			c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
// 			c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
// 			c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
// 			c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT")
// 		}

// 		// Set CORS headers for all requests
// 		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
// 		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
// 		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
// 		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT")

// 		// Handle preflight OPTIONS requests by aborting with status 204
// 		if c.Request.Method == "OPTIONS" {
// 			c.AbortWithStatus(204)
// 			return
// 		}
// 		c.Next()
// 	}
// }

func corsMiddleware() gin.HandlerFunc {
	originsString := "http://localhost:3000"
	allowedOrigins := strings.Split(originsString, ",")

	return func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")

		// Check if the origin is allowed
		if slices.Contains(allowedOrigins, origin) {
			c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
			c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
			c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
			c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT")
		}

		// Always set CORS headers for preflight OPTIONS requests
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		// Proceed to the next middleware for non-OPTIONS requests
		c.Next()
	}
}
