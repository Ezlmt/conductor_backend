package routes

import (
	"conductor_backend/internal/controllers"
	"conductor_backend/internal/middleware"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine) {
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "pong"})
	})
	r.POST("/users/register", controllers.Register)
	r.POST("/users/login", controllers.Login)

	auth := r.Group("/")
	auth.Use(middleware.AuthMiddleware())
	{
		auth.GET("/me", controllers.Me)
	}
}
