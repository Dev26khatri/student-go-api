package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

// Config stores application and database configuration values.
type Config struct {
	AppEnv  string
	AppPort string

	DBHost             string
	DBPort             string
	DBUser             string
	DBPassword         string
	DBName             string
	DBSSLMode          string
	JWTSecret          string
	JWTExpirationHours string
}

// Load reads configuration values from environment variables.
func Load() Config {
	if err := godotenv.Load(); err != nil {
		log.Println("Failed To load Env!.")
	}
	return Config{
		AppEnv:  os.Getenv("APP_ENV"),
		AppPort: os.Getenv("APP_PORT"),

		DBHost:             os.Getenv("DB_HOST"),
		DBPort:             os.Getenv("DB_PORT"),
		DBUser:             os.Getenv("DB_USER"),
		DBPassword:         os.Getenv("DB_PASSWORD"),
		DBName:             os.Getenv("DB_NAME"),
		DBSSLMode:          os.Getenv("DB_SSLMODE"),
		JWTSecret:          os.Getenv("JWT_SECRET"),
		JWTExpirationHours: os.Getenv("JWT_EXPIRATION_HOURS"),
	}
}
