package database

import (
	"conductor_backend/internal/models"
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

const (
	host     = "localhost"
	user     = "user"
	password = "123"
	dbnema   = "conductor"
	port     = 5432
)

func Connect() {
	dsn := fmt.Sprintf("host=%s port=%d user=%s dbname=%s password=%s sslmode=disable", host, port, user, dbnema, password)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	DB = db
	DB.AutoMigrate(&models.User{})
}
