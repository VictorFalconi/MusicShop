package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"server/app/routes"
	"server/config"
)

func init() {
	config.LoadEnvVirables()
	config.ConnectDB()
}

func main() {
	//router := gin.New()
	router := gin.Default()

	routes.UserRouter(router)
	routes.BrandRouter(router)
	routes.ProductRouter(router)
	routes.OrderRouter(router)

	// Use the CORS middleware
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"http://localhost:3000"}
	router.Use(cors.New(config))

	// Start the server
	if err := router.Run(); err != nil {
		panic(err)
	}
}
