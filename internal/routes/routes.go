package routes

import (
	"student-go-service/internal/middleware"
	"student-go-service/internal/student"
	"student-go-service/internal/user"

	"github.com/gin-gonic/gin"
)

// SetupRoutes registers health checks and student API endpoints.
func SetupRoutes(
	router *gin.Engine,
	studentHandler *student.Handler,
	userHandler *user.Handler,
	jwtSecret string,
) {
	router.GET("/helth", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "OK",
		})
	})

	api := router.Group("/api/v1")
	{
		auth := api.Group("/auth")
		{
			auth.POST("/register", userHandler.Register)
			auth.POST("/login", userHandler.Login)
		}

		protected := api.Group("")
		protected.Use(middleware.Auth(jwtSecret))
		students := protected.Group("/students")
		{
			students.POST("", studentHandler.Create)
			students.GET("", studentHandler.GetAll)
			students.GET("/:studentId", studentHandler.GetById)

			students.PUT("/:studentId", middleware.RequireRole("adming"), studentHandler.Update)
			students.DELETE("/:studentId", middleware.RequireRole("admin"), studentHandler.Delete)
		}

	}

}
