package models

import "time"

type Enrollment struct {
	ID uint `gorm:"primaryKey"`
	UserID uint `gorm:"not null"`
	CourseID uint `gorm:"not null"`
	CreatedAt time.Time `gorm:"not null"`
	User User `gorm:"foreignKey:UserID"`
	Course Course `gorm:"foreignKey:CourseID"`
}