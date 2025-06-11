package auth

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
)

const bcryptCost = 10 // Сложность хеширования для bcrypt (10 - стандартное значение, 12 - более безопасно, но медленнее, 14 - очень медленно)

// Hash password func
func HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcryptCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

// ComparePassword func (return error if password does not match)
func ComparePassword(password, hashed string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hashed), []byte(password))
	if err != nil {
		return errors.New("password does not match")
	}
	return nil
}
