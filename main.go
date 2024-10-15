package main

import (
	"api-server/config"
	"api-server/middlewares"
	"api-server/routes"
	"fmt"
	"time"

	"github.com/gin-contrib/cors" // Import the cors package
	"github.com/gin-gonic/gin"
)

func main() {
	// Set Gin to release mode
	gin.SetMode(gin.ReleaseMode)

	// Initialize the database
	config.InitDatabase()

	// Set up the router
	r := gin.Default()

	// Apply CORS middleware
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// Apply the global error handler middleware
	r.Use(middlewares.ErrorHandler())

	// Apply the concurrent request limit middleware
	r.Use(middlewares.LimitConcurrentRequests())

	// Serve static files from the "public/images" directory
	r.Static("/images", "./public/images")
	r.Static("/assets", "./public/assets")

	routes.SetupRoutes(r)

	fmt.Println("Server Running")

	// Start the server
	r.Run(":3000")

}
