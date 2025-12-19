package main

import (
	"conductor_backend/internal/database"
	"conductor_backend/internal/routes"
	"github.com/gin-gonic/gin"
)

func main() {
	database.Connect()
	r := gin.Default()
	routes.RegisterRoutes(r)
	r.Run(":9916")
}
