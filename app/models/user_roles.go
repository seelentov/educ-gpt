package models

type UserRoles struct {
	ID     uint `gorm:"primarykey"`
	RoleID uint `gorm:"not null"`
	Role   Role `gorm:"foreignKey:RoleID"`
	UserID uint `gorm:"not null"`
	User   User `gorm:"foreignKey:UserID"`
}
