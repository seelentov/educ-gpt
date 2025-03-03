package models

type Role struct {
	ID    uint    `gorm:"primarykey"`
	Name  string  `gorm:"unique;not null"`
	Users []*User `gorm:"many2many:user_roles"`
}
