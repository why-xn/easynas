package util

import (
	"encoding/base64"
	"golang.org/x/crypto/bcrypt"
)

// HashPassword generates a bcrypt hash of the password.
func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

// CheckPasswordHash compares a bcrypt hashed password with a plain password.
func CheckPasswordHash(password, hashedPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}

func Base64Encode(data string) string {
	return base64.StdEncoding.EncodeToString([]byte(data))
}

func Base64Decode(data string) string {
	decoded, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		return ""
	}
	return string(decoded)
}
