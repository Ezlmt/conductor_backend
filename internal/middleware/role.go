package middleware

import (
	"conductor_backend/internal/models"

	"github.com/gin-gonic/gin"
)

func RequireProfessor() gin.HandlerFunc {
	return func(c *gin.Context) {
		roleVal, _ := c.Get("role")
		role := roleVal.(int8)
		if role != models.RoleProfessor {
			c.JSON(403, gin.H{"message": "Forbidden", "role": role, "shouldBe": models.RoleProfessor})
			c.Abort()
			return
		}
		c.Next()
	}
}

func RequireStudent() gin.HandlerFunc {
	return func(c *gin.Context) {
		roleVal, _ := c.Get("role")
		role := roleVal.(int8)
		if role != models.RoleStudent {
			c.JSON(403, gin.H{"message": "Forbidden", "role": role, "shouldBe": models.RoleStudent})
			c.Abort()
			return
		}
		c.Next()
	}
}
