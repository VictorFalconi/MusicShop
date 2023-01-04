package routes

import (
	"github.com/gin-gonic/gin"
	"server/app/controller"
)

func UserRouter(router *gin.Engine) {
	router.POST("/register", controller.Register)
	router.POST("/login", controller.Login)
	router.GET("/user", controller.User)

}
