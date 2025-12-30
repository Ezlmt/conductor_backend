package database

import (
	"conductor_backend/internal/models"
	"fmt"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectPostgreSQL() {
	host := getEnv("DB_HOST", "localhost")
	user := getEnv("DB_USER", "user")
	password := getEnv("DB_PASSWORD", "123")
	dbname := getEnv("DB_NAME", "conductor")
	port := getEnv("DB_PORT", "5432")
	dsn := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=disable", host, port, user, dbname, password)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	DB = db
	DB.AutoMigrate(
		&models.User{},
		&models.Course{},
		&models.Enrollment{},
	)
}

func getEnv(key, defaultVal string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return defaultVal
}
