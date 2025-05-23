package middleware

import (
	"net/http"
	"platifyapi/util"

	"github.com/gin-gonic/gin"
)

func Authenticate(c *gin.Context) {
	tokenString := c.GetHeader("Authorization")
	// fmt.Printf("Fetching data")
	userID, err := util.VerifyToken(tokenString)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"message": "failed to verify authentication token",
			"error":   err.Error(),
		})
		return
	}

	c.Set("userID", *userID)
	c.Next()
}
