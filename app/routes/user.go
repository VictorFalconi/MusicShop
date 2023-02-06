package routes

import (
	"github.com/gin-gonic/gin"
	"server/app/controller"
)

func UserRouter(router *gin.Engine) {

	router.POST("/user/register", controller.Register)
	router.POST("/user/login", controller.Login)
	//router.GET("/user/ValidateToken", controller.AuthenticToken)
}
