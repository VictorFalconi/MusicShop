package routes

import (
	"github.com/gin-gonic/gin"
	"server/app/controller"
	"server/middleware"
)

func BrandRouter(router *gin.Engine) {
	router.Use(middleware.CorsMiddleware())
	router.Use(middleware.AuthMiddleware())

	router.POST("/brand", middleware.AdminMiddleware(), controller.CreateBrand)
	router.GET("/brand", controller.ReadBrands)
	router.GET("/brand/:id", controller.ReadBrand)
	router.PUT("/brand/:id", middleware.AdminMiddleware(), controller.UpdateBrand)
	router.DELETE("/brand/:id", middleware.AdminMiddleware(), controller.DeleteBrand)
	router.POST("/brand/file", middleware.AdminMiddleware(), controller.CreateBrand_FromFile)
}
