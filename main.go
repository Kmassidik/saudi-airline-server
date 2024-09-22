package main

import (
	"api-server/config"
	"api-server/middlewares"
	"api-server/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	// Initialize the database
	config.InitDatabase()

	// Set up the router
	r := gin.Default()

	// Apply the global error handler middleware
	r.Use(middlewares.ErrorHandler())

	// Serve static files from the "public/images" directory
	r.Static("/images", "./public/images")
	r.Static("/assets", "./public/assets")

	routes.SetupRoutes(r)

	// Start the server
	r.Run(":3000")
}
