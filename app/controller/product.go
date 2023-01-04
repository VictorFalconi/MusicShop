package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"net/http"
	"server/app/models"
	"server/config"
)

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
