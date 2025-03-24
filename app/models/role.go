package models

type Role struct {
	ID    uint    `gorm:"primarykey" json:"id"`
	Name  string  `gorm:"unique;not null" json:"name"`
	Users []*User `gorm:"many2many:user_roles" json:"-"`
}
