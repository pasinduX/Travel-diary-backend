package integrations

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const (
	AccessTokenTTL  = 15 * time.Minute
	RefreshTokenTTL = 30 * 24 * time.Hour
)

type TokenPair struct {
	AccessToken  string
	RefreshToken string
	RefreshID    string
}

func CreateAccessToken(secret, userID, username, email string) (string, error) {
	claims := jwt.MapClaims{
		"sub":      userID,
		"username": username,
		"email":    email,
		"type":     "access",
		"iat":      time.Now().Unix(),
		"exp":      time.Now().Add(AccessTokenTTL).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

func CreateRefreshToken(secret, userID string) (string, string, error) {
	jti, err := NewTokenID()
	if err != nil {
		return "", "", err
	}
	claims := jwt.MapClaims{
		"sub":  userID,
		"jti":  jti,
		"type": "refresh",
		"iat":  time.Now().Unix(),
		"exp":  time.Now().Add(RefreshTokenTTL).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString([]byte(secret))
	return signed, jti, err
}

func ParseRefreshToken(secret, tokenString string) (*jwt.Token, jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(secret), nil
	})
	if err != nil {
		return nil, nil, err
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, nil, errors.New("invalid refresh token")
	}
	return token, claims, nil
}

func ParseAccessToken(secret, tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(secret), nil
	})
	if err != nil {
		return nil, err
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid access token")
	}
	if tokenType, _ := claims["type"].(string); tokenType != "access" {
		return nil, errors.New("invalid access token")
	}
	return claims, nil
}

func NewTokenID() (string, error) {
	var b [32]byte
	if _, err := rand.Read(b[:]); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(b[:]), nil
}

func HashToken(token string) string {
	sum := sha256.Sum256([]byte(token))
	return hex.EncodeToString(sum[:])
}
