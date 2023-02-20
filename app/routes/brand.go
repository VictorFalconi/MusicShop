package routes

import (
	"github.com/gin-gonic/gin"
	"server/app/controller"
	"server/middleware"
)

func BrandRouter(router *gin.Engine) {
	router.POST("/brand", middleware.AuthMiddleware(), middleware.AdminMiddleware(), controller.CreateBrand)
	router.GET("/brand", controller.ReadBrands)
	router.GET("/brand/:id", controller.ReadBrand)
	router.PUT("/brand/:id", middleware.AuthMiddleware(), middleware.AdminMiddleware(), controller.UpdateBrand)
	router.DELETE("/brand/:id", middleware.AuthMiddleware(), middleware.AdminMiddleware(), controller.DeleteBrand)
	router.POST("/brand/file", middleware.AuthMiddleware(), middleware.AdminMiddleware(), controller.CreateBrand_FromFile)
}
