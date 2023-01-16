package routes

import (
	"github.com/gin-gonic/gin"
	"server/app/controller"
)

func BrandRouter(router *gin.Engine) {
	router.POST("/brand", controller.CreateBrand)
	router.GET("/brand", controller.ReadBrands)
	router.GET("/brand/:id", controller.ReadBrand)
	router.PUT("/brand/:id", controller.UpdateBrand)
	router.DELETE("/brand/:id", controller.DeleteBrand)
	router.POST("/brand/file", controller.CreateBrand_FromFile)
}
