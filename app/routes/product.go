package routes

import (
	"github.com/gin-gonic/gin"
	"server/app/controller"
	"server/middleware"
)

func ProductRouter(router *gin.Engine) {
	router.Use(middleware.Middleware_Authentic())

	router.POST("/product", middleware.Middleware_IsAdmin(), controller.CreateProduct)
	router.GET("/product", controller.ReadProducts)
	router.GET("/product/:id", controller.ReadProduct)
	router.PUT("/product/:id", middleware.Middleware_IsAdmin(), controller.UpdateProduct)
	router.DELETE("/product/:id", middleware.Middleware_IsAdmin(), controller.DeleteProduct)
	//router.POST("/product/file", middleware.Middleware_IsAdmin(), controller.CreateProduct_FromFile)
}

func GalleryRouter(router *gin.Engine) {
	router.Use(middleware.Middleware_Authentic())

	//router.POST("/gallery", middleware.Middleware_IsAdmin(), controller.CreateGallery)
	//router.GET("/gallery", controller.ReadGalleries)
	//router.GET("/gallery/:id", controller.ReadGallery)
	//router.PUT("/gallery/:id", middleware.Middleware_IsAdmin(), controller.UpdateGallery)
	//router.DELETE("/gallery/:id", middleware.Middleware_IsAdmin(), controller.DeleteGallery)
	//router.POST("/gallery/file", middleware.Middleware_IsAdmin(), controller.CreateGallery_FromFile)
}
