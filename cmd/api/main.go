package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/emyu/ecommer-be/config"
	"github.com/emyu/ecommer-be/database"
	"github.com/emyu/ecommer-be/routes"
	"github.com/gin-gonic/gin"
)

func main() {
	// Define flags
	seedFlag := flag.Bool("seed", false, "Seed the database with sample data")
	flag.Parse()

	// Load config
	if err := config.LoadConfig(); err != nil {
		log.Fatal("Failed to load config:", err)
	}

	// Initialize database
	if err := database.InitDB(); err != nil {
		log.Fatal("Failed to initialize database:", err)
	}
	defer database.CloseDB()

	// Seed database if flag is set
	if *seedFlag {
		log.Println("")
		if err := database.SeedDatabase(); err != nil {
			log.Fatal("Failed to seed database:", err)
		}
		return
	}

	// Setup Gin router
	if config.AppConfig.Env == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.Default()

	// Setup routes
	routes.SetupRoutes(router)

	// Start server
	addr := fmt.Sprintf(":%d", config.AppConfig.Port)
	log.Printf("🚀 Server running on http://localhost%s\n", addr)
	log.Printf("📱 API Documentation: Postman collection at /docs\n")
	log.Printf("🔒 Auth: JWT Bearer Token (24 hours expiry)\n")

	if err := router.Run(addr); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
