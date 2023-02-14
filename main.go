package main

import (
	"github.com/gin-gonic/gin"
	"server/app/routes"
	"server/config"
)

func init() {
	config.LoadEnvVirables()
	config.ConnectDB()
}

func main() {
	router := gin.New()

	routes.UserRouter(router)
	routes.BrandRouter(router)
	routes.ProductRouter(router)
	routes.OrderRouter(router)

	router.Run()
}
