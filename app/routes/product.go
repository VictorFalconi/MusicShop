package routes

import (
	"github.com/gin-gonic/gin"
	"server/app/controller"
	"server/middleware"
)

func ProductRouter(router *gin.Engine) {
	router.POST("/product", middleware.AuthMiddleware(), middleware.AdminMiddleware(), controller.CreateProduct)
	router.GET("/product", controller.ReadProducts)
	router.GET("/product/:id", controller.ReadProduct)
	router.PUT("/product/:id", middleware.AuthMiddleware(), middleware.AdminMiddleware(), controller.UpdateProduct)
	router.DELETE("/product/:id", middleware.AuthMiddleware(), middleware.AdminMiddleware(), controller.DeleteProduct)
	router.POST("/product/file", middleware.AuthMiddleware(), middleware.AdminMiddleware(), controller.CreateProduct_FromFile)
}
