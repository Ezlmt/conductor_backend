package controllers

import (
	"conductor_backend/internal/database"
	"conductor_backend/internal/models"
	"time"

	"os"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type registerRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type loginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func Register(c *gin.Context) {
	var req registerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{
			"message": "Invalid request",
		})
		return
	}
	if req.Email == "" || req.Password == "" {
		c.JSON(400, gin.H{
			"message": "Email and password are required",
		})
		return
	}
	hashed, err := bcrypt.GenerateFromPassword(
		[]byte(req.Password),
		bcrypt.DefaultCost,
	)
	if err != nil {
		c.JSON(400, gin.H{
			"message": "Failed to hash password",
		})
		return
	}
	user := models.User{
		Email:        req.Email,
		PasswordHash: string(hashed),
		Role:         models.RoleStudent,
		CreatedAt:    time.Now(),
	}
	if err := database.DB.Create(&user).Error; err != nil {
		c.JSON(400, gin.H{
			"message": "email already exists",
		})
		return
	}
	c.JSON(201, gin.H{
		"id":    user.ID,
		"email": user.Email,
		"role":  user.Role,
	})
}

func Login(c *gin.Context) {
	var req loginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{
			"message": "Invalid request",
		})
		return
	}
	if req.Email == "" {
		c.JSON(400, gin.H{
			"message": "Email is required",
		})
		return
	}
	if req.Password == "" {
		c.JSON(400, gin.H{
			"message": "Password is required",
		})
		return
	}

	user := models.User{}
	if err := database.DB.Where("email = ?", req.Email).First(&user).Error; err != nil {
		c.JSON(401, gin.H{
			"message": "Invalid email or password",
		})
		return
	}
	if bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)) != nil {
		c.JSON(401, gin.H{
			"message": "Invalid email or password",
		})
		return
	}

	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		c.JSON(500, gin.H{
			"message": "Server misconfiguration",
		})
		return
	}

	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":  user.ID,
		"role": user.Role,
		"exp":  time.Now().Add(time.Hour * 24).Unix(),
	}).SignedString([]byte(secret))
	if err != nil {
		c.JSON(500, gin.H{
			"message": "Failed to create token",
		})
		return
	}
	c.JSON(200, gin.H{
		"token": token,
		"user": gin.H{
			"id":    user.ID,
			"email": user.Email,
			"role":  user.Role,
		},
	})

}

func Me(c *gin.Context) {
	userID, ok := c.Get("userID")
	if !ok {
		c.JSON(401, gin.H{"message": "Unauthorized"})
		return
	}
	role, _ := c.Get("role")
	c.JSON(200, gin.H{
		"id":   userID,
		"role": role,
	})
}
