package middleware

import (
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.Method == "OPTIONS" {
			c.Status(204)
			return
		}
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(401, gin.H{
				"message": "Unauthorized",
			})
			c.Abort()
			return
		}
		if !strings.HasPrefix(authHeader, "Bearer ") {
			c.JSON(401, gin.H{
				"message": "Invalid authorization header",
			})
			c.Abort()
			return
		}
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("JWT_SECRET")), nil
		})
		if err != nil || !token.Valid {
			c.JSON(401, gin.H{"message": "Invalid token"})
			c.Abort()
			return
		}
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.JSON(401, gin.H{"message": "Invalid token"})
			c.Abort()
			return
		}
		userID := uint(claims["sub"].(float64))
		role := int8(claims["role"].(float64))

		c.Set("userID", userID)
		c.Set("role", role)
		c.Next()
	}
}
