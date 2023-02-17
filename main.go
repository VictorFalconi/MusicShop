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
	router := gin.Default()

	//// Allow CORS with credentials
	//router.Use(middleware.CorsMiddleware())

	// Use the CORS middleware
	config := cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowCredentials: true,
		ExposeHeaders:    []string{"Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization"},
		AllowMethods:     []string{"POST, GET, OPTIONS, PUT, DELETE"},
		//AllowHeaders:     []string{"Origin"},
		//AllowOriginFunc: func(origin string) bool {return origin == "https://github.com"},
		//MaxAge: 12 * time.Hour,
	}
	router.Use(cors.New(config))

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
