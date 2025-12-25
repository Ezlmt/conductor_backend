package models

import "time"

type Course struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	Name        string    `gorm:"not null" json:"name"`
	Code        string    `gorm:"not null" json:"code"`
	ProfessorID uint      `gorm:"not null" json:"professorId"`
	CreatedAt   time.Time `gorm:"not null" json:"createdAt"`
}
