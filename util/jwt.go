package util

import (
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
)



func GenerateToken(email string, userID int64) (string, error) {
	err := godotenv.Load(".env")
	if err != nil {
		return "", fmt.Errorf("failed to load .env file")
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email": email,
		"user_id": userID,
		"exp":time.Now().Add(time.Minute * 30).Unix(),
	})

	return token.SignedString([]byte(os.Getenv("SECRET_KEY")))
}

func VerifyToken(tokenString string) (*int64, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Verify the signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return []byte(os.Getenv("SECRET_KEY")), nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to parse token")
	}

	tokenIsValid := token.Valid
	if !tokenIsValid {
		return nil, fmt.Errorf("token is not valid")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("unable to parse claims")
	}

	// Check if the token is expired
	if exp, ok := claims["exp"].(float64); ok {
		if time.Unix(int64(exp), 0).Before(time.Now()) {
			return nil, fmt.Errorf("token has expired")
		}
	} else {
		return nil, fmt.Errorf("expiration claim is missing")
	}
	
	// Fetch the user_id from the claims
	userID, ok := claims["user_id"].(float64)
	if !ok {
		return nil, fmt.Errorf("unable to fetch user_id from claims")
	}
	
	convertedUserID := int64(userID)
	return &convertedUserID, nil
}