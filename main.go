package main

import (
	"github.com/gin-gonic/gin"
	"server/app/routes"
	"server/config"
	"server/middleware"
)

func init() {
	config.LoadEnvVirables()
	config.ConnectDB()
}

func main() {
	router := gin.Default()

	// Allow CORS middleware with credentials
	router.Use(middleware.CorsMiddleware())

	// Routes
	routes.UserRouter(router)
	routes.BrandRouter(router)
	routes.ProductRouter(router)
	routes.OrderRouter(router)

	// Start the server
	if err := router.Run(); err != nil {
		panic(err)
	}
}
