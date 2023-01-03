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
		//errors := nil
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
			token, errCreate := auth.Create(currUser.Name)
			if errCreate != nil {
				message := "Internal Server Error"
				errors := errCreate.Error()
				helpers.RespondJSON(ctx, 500, message, errors, nil)
				return
			}
			message := "Login successful!"
			helpers.RespondJSON(ctx, 200, message, nil, token)
			return
		}
	}
}

func ReadUser(ctx *gin.Context) {
	var user models.User
	config.DB.Where("id = ?", ctx.Param("id")).First((&user))
	ctx.JSON(200, &user)
}

func ReadUsers(ctx *gin.Context) {
	var users []models.User
	config.DB.Find(&users)
	ctx.JSON(200, &users)
}

func UpdateUser(ctx *gin.Context) {
	var user models.User
	config.DB.Where("id = ?", ctx.Param("id")).First((&user))
	//ctx.BindJSON(&user)
	//config.DB.Save(&user)
	//ctx.JSON(200, &user)
	err := ctx.ShouldBindJSON(&user)
	if err == nil {
		validate := validator.New()
		err := validate.Struct(&user)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()}) //bad json
		} else {
			err := config.DB.Save(&user).Error
			if err == nil {
				ctx.JSON(200, &user) //oke
			} else {
				ctx.JSON(http.StatusBadGateway, gin.H{"error": err.Error()}) //bad gateway
			}
		}
	}
}

func DeleteUser(ctx *gin.Context) {
	var user models.User
	config.DB.Where("id = ?", ctx.Param("id")).Delete(&user)
	ctx.JSON(200, &user)
}
