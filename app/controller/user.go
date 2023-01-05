package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v4"
	"net/http"
	"server/app/models"
	"server/config"
	"server/helpers"
	"server/middleware"
	"strings"
)

// Register new User
func Register(ctx *gin.Context) {
	var user models.User
	// Check type JSON
	if err := ctx.ShouldBindJSON(&user); err != nil {
		helpers.RespondJSON(ctx, 400, "Error JSON form!", err.Error(), nil)
		return
	}
	// Check validate
	validate := validator.New()
	if err := validate.Struct(&user); err != nil {
		helpers.RespondJSON(ctx, 400, "Error validate!", err.Error(), nil)
		return
	}
	// Hash password & create new User
	user.HashPassword()
	if err := config.DB.Create(&user).Error; err != nil {
		helpers.RespondJSON(ctx, 401, "Duplitace Fields!", err.Error(), nil)
		return
	} else {
		helpers.RespondJSON(ctx, 201, "Created user succesful!", nil, nil)
		return
	}
}

// Login user
func Login(ctx *gin.Context) {
	var currUser models.LoginUser
	if err := ctx.ShouldBindJSON(&currUser); err != nil {
		helpers.RespondJSON(ctx, 400, "Error JSON form!", err.Error(), nil)
		return
	}
	validate := validator.New()
	if err := validate.Struct(&currUser); err != nil {
		helpers.RespondJSON(ctx, 400, "Error validate!", err.Error(), nil)
		return
	}
	// Check Field "name" in db
	user := &models.User{}
	if err := config.DB.Where("name = ?", currUser.Name).First(&user).Error; err != nil {
		helpers.RespondJSON(ctx, 401, "Name doesn't exist", err.Error(), nil)
		return
	} else {
		// Compare password
		if user.ComparePassword(currUser.Password) == false {
			helpers.RespondJSON(ctx, 401, "Incorrect password", err.Error(), nil)
			return
		} else {
			//Create token
			token, errCreate := middleware.CreateToken(user.Id)
			if errCreate != nil {
				helpers.RespondJSON(ctx, 500, "Internal Server Error", errCreate.Error(), nil)
				return
			}
			//Reponse token
			ctx.SetSameSite(http.SameSiteLaxMode)
			ctx.SetCookie("Authorization", token, 3600*12, "", "", false, true)
			helpers.RespondJSON(ctx, 201, "Login successful!", nil, nil)
			return
		}
	}
}

func AuthorizeToken(ctx *gin.Context) {
	//Get token
	bearerToken := ctx.Request.Header.Get("Authorization")
	if len(strings.Split(bearerToken, " ")) == 2 {
		//Validate token
		token, err := middleware.ValidateToken(strings.Split(bearerToken, " ")[1])
		if token.Valid && err == nil {
			claims := token.Claims.(jwt.MapClaims)
			// Query User from id
			var user models.User
			if errDB := config.DB.Where("id = ?", claims["id"]).First(&user).Error; errDB == nil {
				helpers.RespondJSON(ctx, 200, "Validate token successful!", nil, &user)
				return
			} else {
				helpers.RespondJSON(ctx, 400, "User doesn't exist", errDB.Error(), &user)
				return
			}
		} else {
			helpers.RespondJSON(ctx, 401, "Failed to validate token", err.Error(), nil)
			return
		}
	} else {
		helpers.RespondJSON(ctx, 400, "Failed to process request", "No token found", nil)
		return
	}
}
