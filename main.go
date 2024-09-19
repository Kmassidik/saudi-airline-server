package main

import (
	"api-server/config"
	"api-server/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	// Initialize the database
	config.InitDatabase()

	// Set up the router
	r := gin.Default()

	// Serve static files from the "public/images" directory
	r.Static("/images", "./public/images")

	routes.SetupRoutes(r)

	// Start the server
	r.Run(":3000")
}
