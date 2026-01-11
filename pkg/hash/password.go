package hash

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string, salt ...string) (string, error) {
	if password == "" {
		return "", fmt.Errorf("password cannot be empty")
	}

	if len(salt) > 0 {
		password = password + salt[0]
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(hashedPassword), nil
}

func ComparePassword(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}
