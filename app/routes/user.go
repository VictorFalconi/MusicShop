package routes

import (
	"github.com/gin-gonic/gin"
	"server/app/controller"
	"server/middleware"
)

func UserRouter(router *gin.Engine) {
	router.Use(middleware.CorsMiddleware())

	router.POST("/user/register", controller.Register)
	router.GET("/user/login", controller.Login)
	router.PUT("/user", middleware.AuthMiddleware(), controller.UpdateUser) // token
	router.GET("/user", middleware.AuthMiddleware(), controller.ReadUser)   // token
	//router.GET("/user/ValidateToken", controller.AuthenticToken)
}
