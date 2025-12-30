package models

import "time"

const (
	RoleStudent int8 = iota + 1
	RoleProfessor
)

type User struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	Name         string    `gorm:"default:''" json:"name"`
	Email        string    `gorm:"uniqueIndex;not null" json:"email"`
	PasswordHash string    `gorm:"not null" json:"passwordHash"`
	Role         int8      `gorm:"not null" json:"role"`
	CreatedAt    time.Time `gorm:"not null" json:"createdAt"`
}
