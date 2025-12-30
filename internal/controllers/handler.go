package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type CourseHandler struct {
	DB *gorm.DB
}

type UserHandler struct {
	DB  *gorm.DB
	RDB *redis.Client
}

func NewCourseHandler(DB *gorm.DB) *CourseHandler {
	return &CourseHandler{DB: DB}
}

func NewUserHandler(DB *gorm.DB, RDB *redis.Client) *UserHandler {
	return &UserHandler{DB: DB, RDB: RDB}
}

func (h *CourseHandler) Create(c *gin.Context) {
	CreateCourse(c)
}
