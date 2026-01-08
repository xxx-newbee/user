package utils

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)


func EncryptPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("encrypt password failed: %w", err)
	}	
	return string(hashedPassword), nil
}


func ComparePassword(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}