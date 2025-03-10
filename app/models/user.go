package models

import "time"

type User struct {
	ID            uint       `gorm:"primarykey" json:"id"`
	Name          string     `gorm:"unique;not null" json:"name"`
	Email         string     `gorm:"unique;not null" json:"email"`
	Number        string     `gorm:"unique;not null" json:"number"`
	Password      string     `gorm:"not null" json:"-"`
	Roles         []*Role    `gorm:"many2many:user_roles" json:"-"`
	AvatarUrl     string     `json:"avatarUrl"`
	ChatGptModel  string     `gorm:"not null" json:"-"`
	ChatGptToken  string     `gorm:"not null" json:"-"`
	ActivateAt    *time.Time `json:"-"`
	ActivationKey string     `json:"-"`
	CreatedAt     time.Time  `json:"-" sql:"DEFAULT:current_timestamp"`
}
