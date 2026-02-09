package models

type User struct {
	ID uint `gorm:"primaryKey"`
	Name string
	
	RoleID uint
	Role Role
	
	Barangs []Barang `gorm:"many2many:user_barangs; constraint:OnDelete:Cascade"`
}