package router

import (
	"github.com/gin-gonic/gin"
	"server/app/controller"
	"server/app/repository"
	"server/config"
	"server/middleware"
)

func OrderRouter(router *gin.Engine) {
	userRepo := repository.NewUserRepo(config.DB)
	userMiddleware := middleware.NewUserMiddleware(userRepo)

	// User
	router.POST("/order", userMiddleware.AuthMiddleware(), controller.User_CreateOrder)
	router.GET("/order", userMiddleware.AuthMiddleware(), controller.User_ReadOrders)
	router.GET("/order/:id", userMiddleware.AuthMiddleware(), controller.User_ReadOrder)
	router.PUT("/order/:id", userMiddleware.AuthMiddleware(), controller.User_CancelOrder)

	// Admin
	router.GET("/admin_order", userMiddleware.AuthMiddleware(), userMiddleware.AdminMiddleware(), controller.Admin_ReadOrders)
	router.GET("/admin_order/:id", userMiddleware.AuthMiddleware(), userMiddleware.AdminMiddleware(), controller.Admin_ReadOrder)
	router.PUT("/accept_order/:id", userMiddleware.AuthMiddleware(), userMiddleware.AdminMiddleware(), controller.Admin_AcceptOrder) //Notification
	router.PUT("/cancel_order/:id", userMiddleware.AuthMiddleware(), userMiddleware.AdminMiddleware(), controller.Admin_CancelOrder) //Notification

}
