package auth

import (
	"fmt"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
)

var secret []byte

type AuthClaims struct {
	password string
	jwt.StandardClaims
}

func GenerateToken(password string) (string, error) {
	claim := AuthClaims {
		password,
		jwt.StandardClaims {
			Issuer:    "jjungs",
			IssuedAt:  time.Now().Unix(),
			ExpiresAt: time.Now().Add(time.Hour*1 + time.Minute*30).Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
	return token.SignedString(secret)
}

func ValidateToken(tokenString string) (string, error) {
	token, err := jwt.ParseWithClaims(tokenString, &AuthClaims{}, func(t *jwt.Token) (interface{}, error) {
		return secret, nil
	})
	if err != nil {
		return "", err
	}

	if _, ok := token.Claims.(*AuthClaims); ok && token.Valid {
		return "JJUNGS", nil
	} else {
		return "", fmt.Errorf("invalid token")
	}
}

func init() {
	secret = []byte(os.Getenv("SECRET_KEY"))
}