package middleware

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"os"
	"server/app/model"
	"server/app/repository"
	"server/helpers"
	"strings"
	"time"
)

type UserMiddleware struct {
	repo repository.UserRepoInterface
}

func NewUserMiddleware(repo repository.UserRepoInterface) *UserMiddleware {
	return &UserMiddleware{repo}
}

// CORS
func CorsMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		//ctx.Writer.Header().Set("Access-Control-Allow-Origin", "http://localhost:8080")
		//ctx.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		//ctx.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
		//ctx.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		//if ctx.Request.Method == "OPTIONS" {
		//	ctx.AbortWithStatus(204)
		//	return
		//}
		//ctx.Next()

		// Define a list of allowed origins
		allowedOrigins := []string{"http://localhost:3000", "http://localhost:8080"}

		// Middleware that sets the Access-Control-Allow-Origin header
		origin := ctx.Request.Header.Get("Origin")
		for _, allowedOrigin := range allowedOrigins {
			if origin == allowedOrigin {
				ctx.Writer.Header().Set("Access-Control-Allow-Origin", origin)
				ctx.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
				ctx.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
				ctx.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
				if ctx.Request.Method == "OPTIONS" {
					ctx.AbortWithStatus(204)
					return
				}
				break
			}
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

// Authenticaton: Xác thực người dùng
func (um *UserMiddleware) AuthMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		//Get token
		bearerToken := ctx.Request.Header.Get("Authorization")
		if len(strings.Split(bearerToken, " ")) == 2 {
			//Validate token
			token, err := ValidateToken(strings.Split(bearerToken, " ")[1])
			if token.Valid && err == nil { // thiếu nhap random token thi error
				user, errDB := um.repo.GetUserFromToken(token)
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

//Authorization: Ủy quyền người dùng
func (um *UserMiddleware) AdminMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		user := ctx.MustGet("user").(model.User)
		// Find role of account is "admin"

		role, err := um.repo.GetRoleOfUser(&user)
		if err == nil && role.Name == "admin" {
			ctx.Next()
		} else {
			fError := helpers.FieldError{Field: "role", Message: "Account is not authorized, You are not admin"}
			helpers.RespondJSON(ctx, 403, helpers.StatusCodeFromInt(403), fError, nil)
			ctx.Abort()
			return
		}
	}
}
