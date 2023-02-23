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

	router.Use(userMiddleware.AuthMiddleware())

	// User
	router.POST("/order", controller.User_CreateOrder)
	router.GET("/order", controller.User_ReadOrders)
	router.GET("/order/:id", controller.User_ReadOrder)
	router.PUT("/order/:id", controller.User_CancelOrder)

	// Admin
	router.GET("/admin_order", userMiddleware.AdminMiddleware(), controller.Admin_ReadOrders)
	router.GET("/admin_order/:id", userMiddleware.AdminMiddleware(), controller.Admin_ReadOrder)
	router.PUT("/accept_order/:id", userMiddleware.AdminMiddleware(), controller.Admin_AcceptOrder) //Notification
	router.PUT("/cancel_order/:id", userMiddleware.AdminMiddleware(), controller.Admin_CancelOrder) //Notification

}
