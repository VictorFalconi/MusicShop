package router

import (
	"github.com/gin-gonic/gin"
	"server/app/controller"
	"server/app/repository"
	"server/config"
	"server/middleware"
)

func BrandRouter(router *gin.Engine) {
	userRepo := repository.NewUserRepo(config.DB)
	userMiddleware := middleware.NewUserMiddleware(userRepo)

	router.POST("/brand", userMiddleware.AuthMiddleware(), userMiddleware.AdminMiddleware(), controller.CreateBrand)
	router.GET("/brand", controller.ReadBrands)
	router.GET("/brand/:id", controller.ReadBrand)
	router.PUT("/brand/:id", userMiddleware.AuthMiddleware(), userMiddleware.AdminMiddleware(), controller.UpdateBrand)
	router.DELETE("/brand/:id", userMiddleware.AuthMiddleware(), userMiddleware.AdminMiddleware(), controller.DeleteBrand)
	router.POST("/brand/file", userMiddleware.AuthMiddleware(), userMiddleware.AdminMiddleware(), controller.CreateBrand_FromFile)
}
