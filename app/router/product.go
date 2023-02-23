package router

import (
	"github.com/gin-gonic/gin"
	"server/app/controller"
	"server/app/repository"
	"server/config"
	"server/middleware"
)

func ProductRouter(router *gin.Engine) {
	userRepo := repository.NewUserRepo(config.DB)
	userMiddleware := middleware.NewUserMiddleware(userRepo)

	router.POST("/product", userMiddleware.AuthMiddleware(), userMiddleware.AdminMiddleware(), controller.CreateProduct)
	router.GET("/product", controller.ReadProducts)
	router.GET("/product/:id", controller.ReadProduct)
	router.PUT("/product/:id", userMiddleware.AuthMiddleware(), userMiddleware.AdminMiddleware(), controller.UpdateProduct)
	router.DELETE("/product/:id", userMiddleware.AuthMiddleware(), userMiddleware.AdminMiddleware(), controller.DeleteProduct)
	router.POST("/product/file", userMiddleware.AuthMiddleware(), userMiddleware.AdminMiddleware(), controller.CreateProduct_FromFile)
}
