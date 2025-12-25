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
		auth.POST("/courses", middleware.RequireProfessor(), controllers.CreateCourse)
		auth.DELETE("/courses/:id", middleware.RequireProfessor(), controllers.DeleteCourse)
		auth.GET("/courses", controllers.GetCourseByUserID)
		auth.POST("/courses/join", middleware.RequireStudent(), controllers.JoinCourse)
		auth.DELETE("/courses/:id/leave", middleware.RequireStudent(), controllers.LeaveCourse)
		auth.GET("/enrollments", middleware.RequireStudent(), controllers.GetEnrollmentsByStudentID)
		auth.POST("/users/name", controllers.SetName)
	}

	dev := r.Group("/dev")
	dev.Use(middleware.DevOnly())
	{
		dev.DELETE("/courses/:id", controllers.DeleteCourseByID)
		dev.GET("/show-all-courses", controllers.ShowAllCourses)
	}

}
