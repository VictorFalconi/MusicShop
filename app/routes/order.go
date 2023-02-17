package routes

import (
	"github.com/gin-gonic/gin"
	"server/app/controller"
	"server/middleware"
)

func OrderRouter(router *gin.Engine) {
	router.Use(middleware.AuthMiddleware())

	// User
	router.POST("/order", controller.CreateOrder)
	router.GET("/order/:id", controller.ReadOrder)
	router.GET("/order", controller.ReadOrdersOfUser)
	router.PUT("/order/:id", controller.UserCancelOrder)

	// Admin
	router.GET("/orders", middleware.AdminMiddleware(), controller.ReadOrders)
	router.PUT("/accept_order/:id", middleware.AdminMiddleware(), controller.AcceptOrder) //Notification
	router.PUT("/cancel_order/:id", middleware.AdminMiddleware(), controller.CancelOrder) //Notification

}
