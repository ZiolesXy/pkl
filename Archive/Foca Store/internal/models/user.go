package models

type User struct {
	BaseModel
	Name string `gorm:"not null"`
	Email string `gorm:"uinque;not null"`
	Password string `gorm:"not null"`

	RoleID uint
	Role Role
}