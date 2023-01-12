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
	// Check data type
	if err := helpers.DataContentType(ctx, &user); err != nil {
		helpers.RespondJSON(ctx, 400, "Error data type!", err.Error(), nil)
		return
	}
	// Check validate field     ---- Thiếu 1 cái nhập sai, dư field k có trong user
	if err := validator.New().Struct(&user); err != nil {
		listErrors := helpers.ValidateErrors(err.(validator.ValidationErrors))
		helpers.RespondJSON(ctx, 400, "Errors validate!", listErrors, nil)
		return
	}
	// Hash password, create new User (Check validate Database)
	user.HashPassword()
	if err := config.DB.Create(&user).Error; err != nil {
		ErrorDB := helpers.DBError(err)
		helpers.RespondJSON(ctx, 401, "Error Database", ErrorDB, nil)
		return
	} else {
		helpers.RespondJSON(ctx, 201, "Created user successful!", nil, nil)
		return
	}
}

// Login user
func Login(ctx *gin.Context) {
	var currUser models.LoginUser
	if err := helpers.DataContentType(ctx, &currUser); err != nil {
		helpers.RespondJSON(ctx, 400, "Error data type!", err.Error(), nil)
		return
	}
	validate := validator.New()
	if err := validate.Struct(&currUser); err != nil {
		dictErrors := helpers.ValidateErrors(err.(validator.ValidationErrors))
		helpers.RespondJSON(ctx, 400, "Error validate!", dictErrors, nil)
		return
	}
	// Check Field "name" in db
	user := &models.User{}
	if err := config.DB.Where("name = ?", currUser.Name).First(&user).Error; err != nil {
		helpers.RespondJSON(ctx, 401, "Incorrect Filed", "Name isn't already exist", nil)
		return
	} else {
		// Compare password
		if user.ComparePassword(currUser.Password) == false {
			helpers.RespondJSON(ctx, 401, "Incorrect Filed", "Incorrect Password", nil)
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
