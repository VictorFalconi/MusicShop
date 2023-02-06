package routes

import (
	"github.com/gin-gonic/gin"
	"server/app/controller"
	"server/middleware"
)

func BrandRouter(router *gin.Engine) {
	router.Use(middleware.CorsMiddleware())
	router.Use(middleware.Middleware_Authentic())

	router.POST("/brand", middleware.Middleware_IsAdmin(), controller.CreateBrand)
	router.GET("/brand", controller.ReadBrands)
	router.GET("/brand/:id", controller.ReadBrand)
	router.PUT("/brand/:id", middleware.Middleware_IsAdmin(), controller.UpdateBrand)
	router.DELETE("/brand/:id", middleware.Middleware_IsAdmin(), controller.DeleteBrand)
	router.POST("/brand/file", middleware.Middleware_IsAdmin(), controller.CreateBrand_FromFile)
}
