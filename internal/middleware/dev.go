package middleware

import (
	"os"

	"github.com/gin-gonic/gin"
)


func DevOnly() gin.HandlerFunc {
	return func(c *gin.Context) {
		if os.Getenv("APP_ENV") != "dev" {
			c.JSON(403, gin.H{"message": "Forbidden"})
			c.Abort()
			return
		}
		c.Next()
	}
}