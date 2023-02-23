package main

import (
	"github.com/gin-gonic/gin"
	"server/app/router"
	"server/config"
	"server/middleware"
)

func init() {
	config.LoadEnvVirables()
	config.ConnectDB()
}

func main() {
	r := gin.Default()

	// Allow CORS middleware with credentials
	r.Use(middleware.CorsMiddleware())

	// Routers
	router.UserRouter(r)
	router.BrandRouter(r)
	router.ProductRouter(r)
	router.OrderRouter(r)

	// Start the server
	if err := r.Run(); err != nil {
		panic(err)
	}
}
