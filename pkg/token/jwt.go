package token

import (
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"time"
)

func New(userId int64, secret string) (string, error) {
	now := time.Now()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user": userId,
		"nbf":  now.Unix(),
		"exp":  now.Add(15 * time.Minute).Unix(),
		"iat":  now.Unix(),
	})

	return token.SignedString([]byte(secret))
}

func Get(tokenString, secret string) (map[string]interface{}, error) {
	tokenFromString, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(secret), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := tokenFromString.Claims.(jwt.MapClaims); ok {
		return claims, nil
	}

	return nil, fmt.Errorf("invalid token")
}
