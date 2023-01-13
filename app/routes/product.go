package routes

import (
	"github.com/gin-gonic/gin"
	"server/app/controller"
)

func BrandRouter(router *gin.Engine) {
	router.POST("/brand", controller.CreateBrand)
	router.GET("/brand", controller.Register)
	router.GET("/brand/{id}", controller.Register)
	router.PUT("/brand/{id}", controller.Register)
	router.DELETE("/brand/{id}", controller.Register)

}
