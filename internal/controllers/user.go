package controllers

import (
	"conductor_backend/internal/database"
	"conductor_backend/internal/models"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"os"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type registerRequest struct {
	Email    string `json:"email"`
	Name     string `json:"name"`
	Password string `json:"password"`
	Role     int8   `json:"role"`
}

type loginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Role     int8   `json:"role"`
}

func Register(c *gin.Context) {
	var req registerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Println("register error:", err)
		c.JSON(400, gin.H{
			"message": "Invalid request",
		})
		return
	}
	if req.Email == "" || req.Password == "" {
		log.Println("register error: email or password is empty")
		c.JSON(400, gin.H{
			"message": "Email and password are required",
		})
		return
	}
	if req.Role == 0 {
		req.Role = models.RoleStudent
	}
	if req.Name == "" {
		log.Println("register error: name is required")
		c.JSON(400, gin.H{
			"message": "Name is required",
		})
		return
	}
	hashed, err := bcrypt.GenerateFromPassword(
		[]byte(req.Password),
		bcrypt.DefaultCost,
	)
	if err != nil {
		log.Println("register error: failed to hash password")
		c.JSON(400, gin.H{
			"message": "Failed to hash password",
		})
		return
	}
	user := models.User{
		Email:        req.Email,
		Name:         req.Name,
		PasswordHash: string(hashed),
		Role:         req.Role,
		CreatedAt:    time.Now(),
	}
	if err := database.DB.Create(&user).Error; err != nil {
		log.Println("register error: failed to create user")
		c.JSON(400, gin.H{
			"message": "email already exists",
		})
		return
	}
	log.Println("register success: user created")
	c.JSON(201, gin.H{
		"id":    user.ID,
		"email": user.Email,
		"name":  user.Name,
		"role":  user.Role,
	})
}

func Login(c *gin.Context) {
	var req loginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Println("login error: invalid request")
		c.JSON(400, gin.H{
			"message": "Invalid request",
		})
		return
	}
	if req.Email == "" {
		log.Println("login error: email is required")
		c.JSON(400, gin.H{
			"message": "Email is required",
		})
		return
	}
	if req.Password == "" {
		log.Println("login error: password is required")
		c.JSON(400, gin.H{
			"message": "Password is required",
		})
		return
	}

	user := models.User{}
	if err := database.DB.Where("email = ?", req.Email).First(&user).Error; err != nil {
		log.Println("login error: invalid email or password")
		c.JSON(401, gin.H{
			"message": "Invalid email or password",
		})
		return
	}
	if bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)) != nil {
		log.Println("login error: invalid email or password")
		c.JSON(401, gin.H{
			"message": "Invalid email or password",
		})
		return
	}

	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		log.Println("login error: server misconfiguration")
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
		log.Println("login error: failed to create token")
		c.JSON(500, gin.H{
			"message": "Failed to create token",
		})
		return
	}
	log.Println("login success: token created")
	c.JSON(200, gin.H{
		"token": token,
		"user": gin.H{
			"id":    user.ID,
			"name":  user.Name,
			"email": user.Email,
			"role":  user.Role,
		},
	})

}

type setNameRequest struct {
	Name string `json:"name"`
}

func SetName(c *gin.Context) {
	var req setNameRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Println("set name error: invalid request")
		c.JSON(400, gin.H{"message": "Invalid request"})
		return
	}
	if req.Name == "" {
		log.Println("set name error: name is required")
		c.JSON(400, gin.H{"message": "Name is required"})
		return
	}
	userID := c.GetUint("userID")
	if userID == 0 {
		log.Println("set name error: invalid userID")
		c.JSON(401, gin.H{"message": "Unauthorized"})
		return
	}
	// Use Update to only update the Name field, avoiding issues with NULL values
	result := database.DB.Model(&models.User{}).Where("id = ?", userID).Update("name", req.Name)
	if result.Error != nil {
		log.Println("set name error: failed to save user", result.Error)
		c.JSON(500, gin.H{"message": "Failed to save user"})
		return
	}
	if result.RowsAffected == 0 {
		log.Println("set name error: user not found")
		c.JSON(404, gin.H{"message": "User not found"})
		return
	}
	log.Println("set name success: user name set")
	c.JSON(200, gin.H{"message": "User name set successfully"})
}

func Me(c *gin.Context) {
	userID, ok := c.Get("userID")
	key := fmt.Sprintf("user:%d", userID)
	if !ok {
		log.Println("me error: missing userID")
		c.JSON(401, gin.H{"message": "Unauthorized"})
		return
	}
	val, err := database.RDB.Get(database.Ctx, key).Result()
	if err == nil {
		log.Println("me success: user found in cache")
		var user models.User
		json.Unmarshal([]byte(val), &user)
		c.JSON(200, gin.H{
			"id":    user.ID,
			"name":  user.Name,
			"email": user.Email,
			"role":  user.Role,
		})
		return
	}
	log.Println("me success: user not found in cache")
	user := models.User{}
	if err := database.DB.Where("id = ?", userID).First(&user).Error; err != nil {
		log.Println("me error: failed to get user")
		c.JSON(401, gin.H{"message": "Unauthorized"})
		return
	}
	bytes, _ := json.Marshal(user)
	database.RDB.Set(database.Ctx, key, bytes, 5*time.Minute)
	log.Println("me success: userID found")
	c.JSON(200, gin.H{
		"id":    user.ID,
		"name":  user.Name,
		"email": user.Email,
		"role":  user.Role,
	})
}
