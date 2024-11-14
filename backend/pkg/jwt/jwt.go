package jwt

import (
	"fmt"
	"github.com/whyxn/easynas/backend/pkg/db/model"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var jwtKey = []byte("your_secret_key")

type Claims struct {
	User model.User
	jwt.RegisteredClaims
}

// GenerateJWT creates a new JWT token for a given username
func GenerateJWT(user model.User) (string, error) {
	// Set expiration time for token
	expirationTime := time.Now().Add(2 * time.Hour)

	// Create claims with username and expiry
	claims := &Claims{
		User: user,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	// Create the token using the signing method and claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign the token with the secret key
	return token.SignedString(jwtKey)
}

// ValidateJWT parses and validates a JWT token from the request header
func ValidateJWT(tokenStr string) (*Claims, error) {
	claims := &Claims{}

	// Parse and validate the token
	token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		// Validate the signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method")
		}
		return jwtKey, nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	return claims, nil
}
