package main

import (
	"context"
	"log"
	"os"
	"strconv"
	"student-go-service/internal/config"
	"student-go-service/internal/routes"
	"student-go-service/internal/storage"
	"student-go-service/internal/student"
	"student-go-service/internal/user"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {
	// Load application and database settings.
	cfg := config.Load()

	// Connect to the PostgreSQL database.
	db := config.ConnectDatabase(cfg)

	// Create or update the database schema.
	if err := db.AutoMigrate(
		&student.Student{},
		&user.User{},
	); err != nil {
		log.Fatal("Migration Failed", err)
	}

	awsConfig, err := config.LoadAWSConfig()
	if err != nil {
		log.Fatal("Failed to load AWS config 📦:", err)
	}
	s3Storage := storage.NewS3Storage(
		awsConfig,
		os.Getenv("AWS_S3_BUCKET"),
	)
	if err := s3Storage.TestConnection(context.Background()); err != nil {
		log.Fatal("S3 Connection failed:", err)
	}
	log.Println("S3 Connection Successfully 📦")

	// Create the student repository.
	studentRepository := student.NewRepository(db)
	// Create the student service.
	studentService := student.NewService(studentRepository, s3Storage)
	// Create the HTTP handler.
	studentHandler := student.NewHandler(studentService)

	//User Inits
	userRepostitory := user.NewRepository(db)
	userService := user.NewService(userRepostitory)

	hours, err := strconv.Atoi(cfg.JWTExpirationHours)
	if err != nil {
		log.Fatal("Invalid JWT_EXPIRATION Hours")
	}

	userHandler := user.NewHandler(userService, cfg.JWTSecret, time.Duration(hours)*time.Hour)

	// Create the Gin router.
	router := gin.Default()

	// Register all application routes.
	routes.SetupRoutes(
		router,
		studentHandler,
		userHandler,
		cfg.JWTSecret,
	)

	// Start the HTTP server.
	log.Printf("Server Runing on :%s ", cfg.AppPort)
	if err := router.Run(":" + cfg.AppPort); err != nil {
		log.Fatal(err)
	}
}
