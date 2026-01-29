package models

type Role struct {
	ID uint `gorm:"primaryKey"`
	Name string

	Users []User `gorm:"foreignKey:RoleID; constraint:OnDelete:CASCADE;" json:"users,omitempty"`
}