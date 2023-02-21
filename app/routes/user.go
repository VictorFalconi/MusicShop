package routes

import (
	"github.com/gin-gonic/gin"
	"server/app/controller"
	"server/config"
	"server/middleware"
)

func UserRouter(router *gin.Engine) {

	router.POST("/user/register", controller.RegisterHandler(config.DB))
	router.POST("/user/login", controller.Login)
	router.GET("/user", middleware.AuthMiddleware(), controller.ReadUser)   // token
	router.PUT("/user", middleware.AuthMiddleware(), controller.UpdateUser) // token
	//router.GET("/user/ValidateToken", controller.AuthenticToken)
}
