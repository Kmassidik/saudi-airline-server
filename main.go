package main

import (
	"api-server/config"
	"api-server/middlewares"
	"api-server/routes"
	"time"

	"github.com/gin-contrib/cors" // Import the cors package
	"github.com/gin-gonic/gin"
)

func main() {
	// Initialize the database
	config.InitDatabase()

	// Set up the router
	r := gin.Default()

	// Apply CORS middleware
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"},                   // Specify allowed origin
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}, // Allow these HTTP methods
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept"},        // Allow specific headers
		ExposeHeaders:    []string{"Content-Length"},                          // Expose headers if needed
		AllowCredentials: true,                                                // Allow cookies/authentication headers
		MaxAge:           12 * time.Hour,                                      // Cache preflight requests for 12 hours
	}))

	// Apply the global error handler middleware
	r.Use(middlewares.ErrorHandler())

	// Serve static files from the "public/images" directory
	r.Static("/images", "./public/images")
	r.Static("/assets", "./public/assets")

	routes.SetupRoutes(r)

	// Start the server
	r.Run(":3000")
}
