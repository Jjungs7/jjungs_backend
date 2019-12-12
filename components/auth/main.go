package auth

import (
	"fmt"
	"os"

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
			Issuer:    "",
			IssuedAt:  0,
			ExpiresAt: 0,
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

func main() {
	res, err := GenerateToken("IlBdBkm5893")
	fmt.Println("Wow: " + res)
	if err != nil {
		fmt.Println(err)
	}
}