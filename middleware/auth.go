package middleware

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"gorm.io/gorm"
	"os"
	"server/app/models"
	"time"
)

func CreateToken(id uint) (string, error) {
	claims := jwt.MapClaims{}
	claims["authorized"] = true
	claims["id"] = id
	//claims["exp"] = time.Now().Add(time.Hour * 12).Unix() //Token out of date after 12 hours
	claims["exp"] = time.Now().Add(time.Minute * 10).Unix()

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

// Login using username & password
func BasicAuth(currUser models.LoginUser, user models.User, db *gorm.DB, ctx *gin.Context) (string, string, error) {
	if err := db.Where("name = ?", currUser.Name).First(&user).Error; err != nil {
		message := "Incorrect Filed"
		error := "Name isn't already exist"
		return message, error, err
	} else {
		// Compare password
		if user.ComparePassword(currUser.Password) == false {
			message := "Incorrect Filed"
			error := "Incorrect Password"
			return message, error, err
		}
	}
	return "", "", nil
}
