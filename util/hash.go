package util

import (
	"golang.org/x/crypto/bcrypt"
)

// HashPassword takes a plaintext string and returns its bcrypt hash
func HashPassword(plainText string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(plainText), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	
	return string(hashedPassword), nil
}

func CheckPasswordHash(password string, hashedPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}