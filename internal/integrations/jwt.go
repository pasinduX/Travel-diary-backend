package integrations

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func CreateJWT(secret, userID, username, email string) (string, error) {
	claims := jwt.MapClaims{
		"sub":      userID,
		"username": username,
		"email":    email,
		"iat":      time.Now().Unix(),
		"exp":      time.Now().Add(24 * time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}
