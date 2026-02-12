package models

import "time"

type User struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	Name         string    `gorm:"size:100;not null" json:"name"`
	Email        string    `gorm:"uniqueIndex;size:120;not null" json:"email"`
	PasswordHash string    `gorm:"size:255" json:"-"`
	RoleID       uint      `gorm:"not null" json:"role_id"`
	Role         Role      `gorm:"constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;" json:"role"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}
