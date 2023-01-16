package middleware

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"gorm.io/gorm"
	"os"
	"server/app/models"
	"server/config"
	"server/helpers"
	"strings"
	"time"
)

func CreateToken(id uint) (string, error) {
	claims := jwt.MapClaims{}
	claims["authorized"] = true
	claims["id"] = id
	claims["exp"] = time.Now().Add(time.Minute * 30).Unix()

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(os.Getenv("SECRET_JWT")))
}

func ValidateToken(token string) (*jwt.Token, error) {
	return jwt.Parse(token, func(t_ *jwt.Token) (interface{}, error) {
		if _, ok := t_.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method %v", t_.Header["alg"])
		}
		return []byte(os.Getenv("SECRET_JWT")), nil
	})
}

// Authenticaton: Xác thực người dùng
func Middleware_Authentic() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		//Get token
		bearerToken := ctx.Request.Header.Get("Authorization")
		if len(strings.Split(bearerToken, " ")) == 2 {
			//Validate token
			token, err := ValidateToken(strings.Split(bearerToken, " ")[1])
			if token.Valid && err == nil {
				claims := token.Claims.(jwt.MapClaims)
				// Query User from id
				var user models.User
				if errDB := config.DB.Where("id = ?", claims["id"]).First(&user).Error; errDB == nil {
					ctx.Next()
				} else {
					helpers.RespondJSON(ctx, 400, "User doesn't exist", errDB.Error(), &user)
					ctx.Abort()
					return
				}
			} else {
				helpers.RespondJSON(ctx, 401, "Failed to validate token", err.Error(), nil)
				ctx.Abort()
				return
			}
		} else {
			helpers.RespondJSON(ctx, 400, "Failed to process request", "No token found", nil)
			ctx.Abort()
			return
		}
	}
}

// Authorization: Ủy quyền người dùng
func Middleware_IsAdmin() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		//Get token
		bearerToken := ctx.Request.Header.Get("Authorization")
		if len(strings.Split(bearerToken, " ")) == 2 {
			//Validate token
			token, err := ValidateToken(strings.Split(bearerToken, " ")[1])
			if token.Valid && err == nil {
				claims := token.Claims.(jwt.MapClaims)
				// Query User from id
				var user models.User
				if errDB := config.DB.Where("id = ?", claims["id"]).First(&user).Error; errDB == nil {
					// Find role of account is "admin"
					errName, RoleName := user.GetNameRoleUser(config.DB)
					if errName == nil && RoleName == "admin" {
						ctx.Next()
					} else {
						helpers.RespondJSON(ctx, 400, "Error Authorization", "Account is not authorized, You are not admin", nil)
						ctx.Abort()
						return
					}
				}
			}
		}
	}
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
