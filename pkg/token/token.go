package token

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func GenerateToken(username string, secret string) (string, error) {
	claims := jwt.MapClaims{
		"sub": username,                              // кому выдан
		"exp": time.Now().Add(24 * time.Hour).Unix(), // когда истекает
	}
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := t.SignedString([]byte(secret)) // signed — готовая строка токена
	return signed, err
}

func ParseToken(token string, secret string) (string, error) {
	// 1. разобрать и проверить токен
	parsed, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	}, jwt.WithValidMethods([]string{"HS256"}))
	if err != nil {
		return "", err
	}

	// 3. достать claims (приведение типа)
	claims, ok := parsed.Claims.(jwt.MapClaims)
	if !ok {
		return "", errors.New("invalid claims")
	}

	// затем claims["sub"] — это interface{}, его тоже нужно привести к string
	sub, ok := claims["sub"].(string)
	if !ok {
		return "", errors.New("invalid sub")
	}
	return sub, nil

}
