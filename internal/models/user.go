package models

import "time"

const (
	RoleStudent int8 = iota + 1
	RoleProfessor
)

type User struct {
	ID           uint   `gorm:"primaryKey"`
	Email        string `gorm:"uniqueIndex;not null"`
	PasswordHash string `gorm:"not null"`
	Role         int8   `gorm:"not null"`
	CreatedAt    time.Time
}
