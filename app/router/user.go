package router

import (
	"github.com/gin-gonic/gin"
	"server/app/controller"
	"server/app/repository"
	"server/app/service"
	"server/config"
	"server/middleware"
)

func UserRouter(router *gin.Engine) {
	userRepo := repository.NewUserRepo(config.DB)
	userService := service.NewUserService(userRepo)
	userMiddleware := middleware.NewUserMiddleware(userRepo)
	userController := controller.NewUserController(userService)

	router.POST("/user/register", userController.RegisterHandler())
	router.POST("/user/login", userController.LoginHandler())
	// OAuth
	router.GET("/auth", userController.OAuth2Home())
	router.GET("/auth/login", userController.OAuth2LoginHandler())
	router.GET("/auth/callback", userController.OAuth2CallbackHandler())

	router.GET("/user", userMiddleware.AuthMiddleware(), userController.ReadUserHandler())
	router.PUT("/user", userMiddleware.AuthMiddleware(), userController.UpdateUserHandler())

}
