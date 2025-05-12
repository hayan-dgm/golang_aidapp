package main

import (
	"aidapp_api_golang/db"
	"aidapp_api_golang/handlers"
	"aidapp_api_golang/middleware"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	// Initialize database
	db.InitDB()
	defer db.CloseDB()

	// Pass DB connection to middleware
	middleware.Initialize(db.DB)

	// Start WebSocket hub
	go handlers.HubInstance.Run()

	// Create router
	r := gin.Default()

	// Routes
	r.GET("/", handlers.Home)
	r.POST("/login", handlers.Login)
	r.GET("/ws", handlers.HandleWebSocket)

	// Authenticated routes
	auth := r.Group("/")
	auth.Use(middleware.JWTMiddleware())
	{
		auth.POST("/logout", handlers.Logout)
		auth.GET("/families", handlers.GetFamilies)
		auth.GET("/families/:id", handlers.GetFamily)
		auth.PUT("/families/:id/products", handlers.UpdateProducts)
		auth.POST("/families", handlers.AddFamily)
		auth.GET("/active_sessions", handlers.GetActiveSessions)
		auth.DELETE("/active_sessions", handlers.ClearActiveSessions)
	}

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("Server running on port %s", port)
	r.Run(":" + port)
}
