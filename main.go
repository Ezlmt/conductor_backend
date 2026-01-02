package main

import (
	"conductor_backend/internal/database"
	"conductor_backend/internal/routes"

	"log"
	"os"
	"strings"

	_ "conductor_backend/docs"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

// @title Conductor API
// @version 1.0
// @description Conductor backend API
// @termsOfService https://example.com/terms/

// @contact.name API Support
// @contact.email support@example.com

// @license.name MIT

// @openapi 3.0
func main() {
	_ = godotenv.Load()

	log.SetFlags(log.LstdFlags | log.Lshortfile)
	log.SetOutput(os.Stdout)
	log.Println("Starting server...")
	database.ConnectPostgreSQL()
	database.ConnectRedis()
	r := gin.Default()
	r.Use(cors.New(cors.Config{
		AllowOrigins:     getCorsOrigins(),
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		AllowCredentials: true,
	}))
	routes.RegisterRoutes(r)
	r.Run(":9916")
}

func getCorsOrigins() []string {
	origins := os.Getenv("CORS_ALLOWED_ORIGINS")
	if origins == "" {
		return []string{}
	}
	return strings.Split(origins, ",")
}
