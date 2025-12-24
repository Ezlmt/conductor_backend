package controllers

import (
	"conductor_backend/internal/database"
	"conductor_backend/internal/models"
	"errors"
	"log"
	"time"

	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type createCourseRequest struct {
	Name string `json:"name"`
	Code string `json:"code"`
}

func CreateCourse(c *gin.Context) {
	var req createCourseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Println("create course error: invalid request")
		c.JSON(400, gin.H{"message": "Invalid request"})
		return
	}
	course := models.Course{}
	err := database.DB.Where("name = ?", req.Name).First(&course).Error
	if err == nil {
		log.Println("create course error: course already exists")
		c.JSON(400, gin.H{"message": "Course already exists"})
		return
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		log.Println("create course error: database error")
		c.JSON(500, gin.H{"message": "database error"})
		return
	}
	course = models.Course{
		Name:        req.Name,
		Code:        req.Code,
		ProfessorID: c.GetUint("userID"),
		CreatedAt:   time.Now(),
	}
	if err := database.DB.Create(&course).Error; err != nil {
		log.Println("create course error: failed to create course")
		c.JSON(400, gin.H{"message": "Failed to create course"})
		return
	}
	log.Println("create course success: course created")
	c.JSON(201, gin.H{
		"id":          course.ID,
		"name":        course.Name,
		"code":        course.Code,
		"professorID": course.ProfessorID,
		"createdAt":   course.CreatedAt,
	})
}

type deleteCourseRequest struct {
	Id uint `json:"id"`
}

func DeleteCourse(c *gin.Context) {
	var req deleteCourseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Println("delete course error: invalid request")
		c.JSON(400, gin.H{"message": "Invalid request"})
		return
	}
	result := database.DB.Delete(&models.Course{}, req.Id)
	if result.Error != nil {
		log.Println("delete course error: failed to delete course")
		c.JSON(500, gin.H{"message": "Failed to delete course"})
		return
	}
	if result.RowsAffected == 0 {
		log.Println("delete course error: course not found")
		c.JSON(404, gin.H{"message": "Course not found"})
		return
	}
	log.Println("delete course success: course deleted")
	c.JSON(200, gin.H{"message": "Course deleted successfully"})
}

func GetCourseByUserID(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		log.Println("get course by userID error: missing userID")
		c.JSON(401, gin.H{"message": "Missing userID"})
		return
	}
	courses := []models.Course{}
	if err := database.DB.Where("professor_id = ?", userID).Find(&courses).Error; err != nil {
		log.Println("get course by userID error: failed to get courses")
		c.JSON(400, gin.H{"message": "Failed to get courses"})
		return
	}
	log.Println("get course by userID success: courses found")
	c.JSON(200, gin.H{"courses": courses})
}

type joinCourseRequest struct {
	Code string `json:"code"`
}

func JoinCourse(c *gin.Context) {
	var req joinCourseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Println("join course error: invalid request")
		c.JSON(400, gin.H{"message": "Invalid request"})
		return
	}

	course := models.Course{}
	if err := database.DB.Where("code=?", req.Code).First(&course).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Println("join course error: course not found")
			c.JSON(404, gin.H{"message": "Course not found"})
			return
		}
		log.Println("join course error: failed to get course")
		c.JSON(500, gin.H{"message": "Failed to get course"})
		return
	}
	log.Println("join course success: course found")
	existingEnrollment := models.Enrollment{}
	err := database.DB.Where("course_id = ? AND user_id = ?", course.ID, c.GetUint("userID")).First(&existingEnrollment).Error
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			log.Println("join course error: failed to check existing enrollment")
			c.JSON(500, gin.H{"message": "Failed to check existing enrollment"})
			return
		}
	}
	if existingEnrollment.ID != 0 {
		log.Println("join course error: already enrolled in this course")
		c.JSON(400, gin.H{"message": "Already enrolled in this course"})
		return
	}
	enrollment := models.Enrollment{
		UserID:    c.GetUint("userID"),
		CourseID:  course.ID,
		CreatedAt: time.Now(),
	}
	if err := database.DB.Create(&enrollment).Error; err != nil {
		log.Println("join course error: failed to join course")
		c.JSON(500, gin.H{"message": "Failed to join course"})
		return
	}
	log.Println("join course success: joined course")
	c.JSON(200, gin.H{
		"message":  "Joined course successfully",
		"courseId": course.ID})
}

func LeaveCourse(c *gin.Context) {
	courseID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		log.Println("leave course error: invalid course ID")
		c.JSON(400, gin.H{"message": "Invalid course ID"})
		return
	}
	userID := c.GetUint("userID")
	var enrollment models.Enrollment
	err = database.DB.Where("course_id = ? AND user_id = ?", courseID, userID).First(&enrollment).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Println("leave course error: enrollment not found")
			c.JSON(400, gin.H{"message": "Enrollment not found"})
			return
		}
		log.Println("leave course error: failed to get enrollment")
		c.JSON(500, gin.H{"message": "Failed to get enrollment"})
		return
	}
	result := database.DB.Delete(&enrollment)
	if result.Error != nil {
		log.Println("leave course error: failed to leave course")
		c.JSON(500, gin.H{"message": "Failed to leave course"})
		return
	}
	log.Println("leave course success: unenrolled course")
	c.JSON(200, gin.H{"message": "Unenrolled course successfully"})
}

type enrollmentResponse struct {
	CouresId   uint      `json:"courseId"`
	CouresName string    `json:"courseName"`
	CouresCode string    `json:"courseCode"`
	JoinedAt   time.Time `json:"joinedAt"`
}

func GetEnrollmentsByStudentID(c *gin.Context) {
	studentID, exists := c.Get("userID")
	if !exists {
		log.Println("get enrollments by studentID error: missing userID")
		c.JSON(401, gin.H{"message": "Missing userID"})
		return
	}
	enrollments := []models.Enrollment{}
	err := database.DB.Where("user_id = ?", studentID).Preload("Course").Find(&enrollments).Error
	if err != nil {
		log.Println("get enrollments by studentID error: failed to get enrollments")
		c.JSON(500, gin.H{"message": "Failed to get enrollments"})
		return
	}
	log.Println("get enrollments by studentID success: enrollments found")
	responses := []enrollmentResponse{}
	for e := range enrollments {
		responses = append(responses, enrollmentResponse{
			CouresId:   enrollments[e].Course.ID,
			CouresName: enrollments[e].Course.Name,
			CouresCode: enrollments[e].Course.Code,
			JoinedAt:   enrollments[e].CreatedAt,
		})
	}
	log.Println("get enrollments by studentID success: responses found")
	c.JSON(200, gin.H{"courses": responses})
}

// ------------------------------------------------------
// All functions below are only for development purposes
// ------------------------------------------------------
func DeleteCourseByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		log.Println("delete course by ID error: invalid course ID")
		c.JSON(400, gin.H{"message": "Invalid course ID"})
		return
	}
	result := database.DB.Delete(&models.Course{}, id)
	if result.Error != nil {
		log.Println("delete course by ID error: failed to delete course")
		c.JSON(500, gin.H{"message": "Failed to delete course"})
		return
	}
	if result.RowsAffected == 0 {
		log.Println("delete course by ID error: course not found")
		c.JSON(404, gin.H{"message": "Course not found"})
		return
	}
	log.Println("delete course by ID success: course deleted")
	c.JSON(200, gin.H{"message": "Course deleted successfully"})
}

func ShowAllCourses(c *gin.Context) {
	courses := []models.Course{}
	if err := database.DB.Find(&courses).Error; err != nil {
		log.Println("show all courses error: failed to get courses")
		c.JSON(400, gin.H{"message": "Failed to get courses"})
		return
	}
	log.Println("show all courses success: courses found")
	c.JSON(200, gin.H{"courses": courses})
}
