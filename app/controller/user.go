package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"net/http"
	"server/app/models"
	"server/config"
	"server/helpers"
	"server/middleware"
)

// Register new User
func Register(ctx *gin.Context) {
	var user models.User
	// Check type JSON
	if err := ctx.ShouldBindJSON(&user); err != nil {
		message := "Error JSON form!"
		errors := err.Error()
		helpers.RespondJSON(ctx, 400, message, errors, nil)
		return
	}
	// Check validate
	validate := validator.New()
	if err := validate.Struct(&user); err != nil {
		message := "Error validate!"
		errors := err.Error()
		helpers.RespondJSON(ctx, 400, message, errors, nil)
		return
	}
	// Hash password & create new User
	user.HashPassword()
	if err := config.DB.Create(&user).Error; err != nil {
		message := "Duplitace Fields!"
		errors := err.Error()
		helpers.RespondJSON(ctx, 401, message, errors, nil)
		return
	} else {
		message := "Created user succesful!"
		helpers.RespondJSON(ctx, 201, message, nil, nil)
		return
	}
}

// Login user
func Login(ctx *gin.Context) {
	var currUser models.LoginUser
	if err := ctx.ShouldBindJSON(&currUser); err != nil {
		message := "Error JSON form!"
		errors := err.Error()
		helpers.RespondJSON(ctx, 400, message, errors, nil)
		return
	}
	validate := validator.New()
	if err := validate.Struct(&currUser); err != nil {
		message := "Error validate!"
		errors := err.Error()
		helpers.RespondJSON(ctx, 400, message, errors, nil)
		return
	}
	// Check Field "name" in db
	user := &models.User{}
	if err := config.DB.Where("name = ?", currUser.Name).First(&user).Error; err != nil {
		message := "Name doesn't exist"
		errors := err.Error()
		helpers.RespondJSON(ctx, 401, message, errors, nil)
		return
	} else {
		// Compare password
		if user.ComparePassword(currUser.Password) == false {
			message := "Incorrect password"
			errors := err.Error()
			helpers.RespondJSON(ctx, 401, message, errors, nil)
			return
		} else {
			//Create token
			token, errCreate := middleware.CreateToken(currUser.Name)
			if errCreate != nil {
				message := "Internal Server Error"
				errors := errCreate.Error()
				helpers.RespondJSON(ctx, 500, message, errors, nil)
				return
			}
			//Reponse
			ctx.SetSameSite(http.SameSiteLaxMode)
			ctx.SetCookie("Authorization", token, 3600*12, "", "", false, true)
			message := "Login successful!"
			helpers.RespondJSON(ctx, 201, message, nil, nil)
			return
		}
	}
}

func getNameUser(ctx *gin.Context) {
	cookie, err := ctx.Cookie("Authorization")
	if err == nil {
		message := "Dont get Cookie"
		errors := err.Error()
		helpers.RespondJSON(ctx, 404, message, errors, nil)
		return
	} else {
		message := "its cookie"
		helpers.RespondJSON(ctx, 200, message, nil, cookie)
		return
	}
}
