package routes

import (
	"github.com/gin-gonic/gin"
	"server/app/controller"
	"server/middleware"
)

func ProductRouter(router *gin.Engine) {
	router.Use(middleware.AuthMiddleware())

	router.POST("/product", middleware.AdminMiddleware(), controller.CreateProduct)
	router.GET("/product", controller.ReadProducts)
	router.GET("/product/:id", controller.ReadProduct)
	router.PUT("/product/:id", middleware.AdminMiddleware(), controller.UpdateProduct)
	router.DELETE("/product/:id", middleware.AdminMiddleware(), controller.DeleteProduct)
	router.POST("/product/file", middleware.AdminMiddleware(), controller.CreateProduct_FromFile)
}
