package routes

import (
	"github.com/gin-gonic/gin"
	"server/app/controller"
	"server/middleware"
)

func OrderRouter(router *gin.Engine) {
	router.Use(middleware.AuthMiddleware())

	// User
	router.POST("/order", controller.User_CreateOrder)
	router.GET("/order", controller.User_ReadOrders)
	router.GET("/order/:id", controller.User_ReadOrder)
	router.PUT("/order/:id", controller.User_CancelOrder)

	// Admin
	router.GET("/orders", middleware.AdminMiddleware(), controller.Admin_ReadOrders)
	router.GET("/order/:id", middleware.AdminMiddleware(), controller.Admin_ReadOrder)
	router.PUT("/accept_order/:id", middleware.AdminMiddleware(), controller.Admin_AcceptOrder) //Notification
	router.PUT("/cancel_order/:id", middleware.AdminMiddleware(), controller.Admin_CancelOrder) //Notification

}
