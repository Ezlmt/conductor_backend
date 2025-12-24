package models

import "time"

type Course struct {
	ID uint `gorm:"primaryKey"`
	Name string `gorm:"not null"`
	Code string `gorm:"not null"`
	ProfessorID uint `gorm:"not null"`
	CreatedAt time.Time `gorm:"not null"`
}