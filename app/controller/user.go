package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"net/http"
	"server/app/models"
	"server/config"
	"server/helpers"
)

// Register new User
func Register(ctx *gin.Context) {
	var user models.User
	// Check type JSON
	if err := ctx.ShouldBindJSON(&user); err != nil {
		message := err.Error()
		helpers.RespondJSON(ctx, 400, message, nil)
		return
	}
	// Check validate
	validate := validator.New()
	if err := validate.Struct(&user); err != nil {
		message := err.Error()
		helpers.RespondJSON(ctx, 401, message, nil)
		return
	}
	// Hash password & create new User
	user.HashPassword()
	if err := config.DB.Create(&user).Error; err != nil {
		message := err.Error()
		helpers.RespondJSON(ctx, 401, message, nil)
		return
	} else {
		message := "Created user succesful!"
		helpers.RespondJSON(ctx, 201, message, nil)
		return
	}
}

// Login user
func Login(ctx *gin.Context) {
	var currUser models.LoginUser
	if err := ctx.ShouldBindJSON(&currUser); err != nil {
		message := err.Error()
		helpers.RespondJSON(ctx, 400, message, nil)
		return
	}
	validate := validator.New()
	if err := validate.Struct(&currUser); err != nil {
		message := err.Error()
		helpers.RespondJSON(ctx, 401, message, nil)
		return
	}
	// Check Field "name" in db
	user := &models.User{}
	if err := config.DB.Where("name = ?", currUser.Name).First(&user).Error; err != nil {
		message := "Name doesn't exist"
		helpers.RespondJSON(ctx, 401, message, nil)
		return
	} else {
		// Compare password
		if user.ComparePassword(currUser.Password) == false {
			message := "Incorrect password"
			helpers.RespondJSON(ctx, 401, message, nil)
			return
		} else {
			message := "Login successful!"
			helpers.RespondJSON(ctx, 200, message, nil)
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
