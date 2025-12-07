package services

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const (
	expTime = 8
	Key     = "1234567891234567"
)

var JwtKey = []byte("my_secret_key")

type Claims struct {
	Hesh [32]byte `json:"hesh"`
	jwt.RegisteredClaims
}

func GenerateJWT(data [32]byte) (string, error) {
	experationTime := time.Now().Add(expTime * time.Hour)
	claims := &Claims{
		Hesh: data,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(experationTime),
			Issuer:    "Planner",
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(JwtKey)
}

func ParseToken(tokenString string) (*Claims, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		// Проверяем, что метод подписи соответствует ожидаемому
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(JwtKey), nil // Возвращаем секретный ключ
	})

	if err != nil {
		return nil, err
	}
	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}
	return claims, nil
}
