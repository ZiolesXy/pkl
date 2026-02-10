package model

import "gorm.io/gorm"

const (
	RoleAdmin    = "admin"
	RoleCustomer = "customer"
)

type User struct {
	gorm.Model
	Name     string `gorm:"type:varchar(255);not null" json:"name"`
	Email    string `gorm:"type:varchar(255);unique;not null" json:"email"`
	Password string `gorm:"type:varchar(255);not null" json:"-"`
	Role     string `gorm:"type:varchar(50);default:'customer'" json:"role"`
}

type RegisterInput struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

type LoginInput struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}