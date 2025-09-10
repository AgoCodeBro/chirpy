package auth

import (
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	passwordByte := []byte(password)
	hashedPasswordBytes, err := bcrypt.GenerateFromPassword(passwordByte, 4)
	if err != nil {
		return "", err
	}

	hashedPassword := string(hashedPasswordBytes)

	return hashedPassword, nil
}

func CheckPasswordHash(password string, hash string) error {
	passwordByte := []byte(password)
	hashByte := []byte(hash)

	return bcrypt.CompareHashAndPassword(hashByte, passwordByte)
}
