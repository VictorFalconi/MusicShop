package middleware

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"os"
	"server/app/models"
	"server/config"
	"server/helpers"
	"strings"
	"time"
)

// CORS
func CorsMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Writer.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
		ctx.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		ctx.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
		ctx.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		if ctx.Request.Method == "OPTIONS" {
			ctx.AbortWithStatus(204)
			return
		}
		ctx.Next()
	}
}

func CreateToken(id uint) (string, error) {
	claims := jwt.MapClaims{}
	claims["authorized"] = true
	claims["id"] = id
	claims["exp"] = time.Now().Add(time.Hour * 1).Unix()

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

func GetUserFromToken(token *jwt.Token) (interface{}, error) {
	claims := token.Claims.(jwt.MapClaims)
	// Query User from id
	var user models.User
	if err := config.DB.Where("id = ?", claims["id"]).First(&user).Error; err != nil {
		return nil, err
	} else {
		return user, nil
	}
}

// Authenticaton: Xác thực người dùng
func AuthMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		//Get token
		bearerToken := ctx.Request.Header.Get("Authorization")
		if len(strings.Split(bearerToken, " ")) == 2 {
			//Validate token
			token, err := ValidateToken(strings.Split(bearerToken, " ")[1])
			if token.Valid && err == nil { // thiếu nhap random token thi error
				user, errDB := GetUserFromToken(token)
				if errDB == nil {
					ctx.Set("user", user)
					ctx.Next()
				} else {
					status, ErrorDB := helpers.DBError(errDB)
					helpers.RespondJSON(ctx, status, helpers.StatusCodeFromInt(status), ErrorDB, nil)
					ctx.Abort()
					return
				}
			} else {
				fError := helpers.FieldError{Field: "token", Message: err.Error()}
				helpers.RespondJSON(ctx, 401, helpers.StatusCodeFromInt(401), fError, nil)
				ctx.Abort()
				return
			}
		} else {
			fError := helpers.FieldError{Field: "token", Message: "No token found"}
			helpers.RespondJSON(ctx, 400, helpers.StatusCodeFromInt(400), fError, nil)
			ctx.Abort()
			return
		}
	}
}

// Authorization: Ủy quyền người dùng
func AdminMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		user := ctx.MustGet("user").(models.User)
		// Find role of account is "admin"
		errName, RoleName := user.GetNameRoleUser(config.DB)
		if errName == nil && RoleName == "admin" {
			ctx.Next()
		} else {
			fError := helpers.FieldError{Field: "role", Message: "Account is not authorized, You are not admin"}
			helpers.RespondJSON(ctx, 403, helpers.StatusCodeFromInt(403), fError, nil)
			ctx.Abort()
			return
		}
	}
}

//// Login using username & password
//func BasicAuth(currUser models.LoginUser, user models.User, db *gorm.DB, ctx *gin.Context) (string, string, error) {
//	if err := db.Where("name = ?", currUser.Name).First(&user).Error; err != nil {
//		message := "Incorrect Filed"
//		error := "Name isn't already exist"
//		return message, error, err
//	} else {
//		// Compare password
//		if user.ComparePassword(currUser.Password) == false {
//			message := "Incorrect Filed"
//			error := "Incorrect Password"
//			return message, error, err
//		}
//	}
//	return "", "", nil
//}
