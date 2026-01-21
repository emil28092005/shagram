package auth

import (
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
}

func secret() ([]byte, error) {
	s := os.Getenv("JWT_SECRET")
	if s == "" {
		return nil, errors.New("JWT_SECRET is not set")
	}
	return []byte(s), nil
}

func NewAccessToken(username string, ttl time.Duration) (string, error) {
	key, err := secret()
	if err != nil {
		return "", err
	}
	now := time.Now()
	claims := Claims{
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   username,
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(ttl)),
		},
	}

	t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return t.SignedString(key)
}

func ParseAccessToken(tokenString string) (*Claims, error) {
	key, err := secret()
	if err != nil {
		return nil, err
	}

	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (any, error) {
		return key, nil
	},
		jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}),
	)
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}
