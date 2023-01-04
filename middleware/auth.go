package middleware

import (
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"os"
	"time"
)

func CreateToken(id uint) (string, error) {
	claims := jwt.MapClaims{}
	claims["authorized"] = true
	claims["id"] = id
	//claims["exp"] = time.Now().Add(time.Hour * 12).Unix() //Token out of date after 12 hours
	claims["exp"] = time.Now().Add(time.Second * 30).Unix()

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(os.Getenv("SECRET_JWT")))
}

func ValidateToken(token string) (*jwt.Token, error) {
	return jwt.Parse(token, func(t_ *jwt.Token) (interface{}, error) {
		if _, ok := t_.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method %v", t_.Header["alg"])
		}
		return []byte(os.Getenv("SECRET_JWT")), nil
	})
}
