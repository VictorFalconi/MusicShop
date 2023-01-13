package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"net/http"
	"server/app/models"
	"server/config"
	"server/helpers"
)

func CreateBrand(ctx *gin.Context) {
	var brand models.Brand
	// Check data type
	if err := helpers.DataContentType(ctx, &brand); err != nil {
		helpers.RespondJSON(ctx, 400, "Error data type!", err.Error(), nil)
		return
	}
	// Check validate field
	if err := validator.New().Struct(&brand); err != nil {
		listErrors := helpers.ValidateErrors(err.(validator.ValidationErrors))
		helpers.RespondJSON(ctx, 400, "Errors validate!", listErrors, nil)
		return
	}
	// Create new Brand (Check validate Database)
	if err := config.DB.Create(&brand).Error; err != nil {
		ErrorDB := helpers.DBError(err)
		helpers.RespondJSON(ctx, 401, "Error Database", ErrorDB, nil)
		return
	} else {
		helpers.RespondJSON(ctx, 201, "Created brand successful!", nil, nil)
		return
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
