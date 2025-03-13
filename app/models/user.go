package models

import "time"

type User struct {
	ID            uint       `gorm:"primarykey" json:"id"`
	Name          string     `gorm:"unique;not null" json:"name"`
	Email         string     `gorm:"unique;not null" json:"email"`
	Password      string     `gorm:"not null" json:"-"`
	Roles         []*Role    `gorm:"many2many:user_roles" json:"-"`
	AvatarUrl     string     `json:"avatar_url"`
	ChatGptModel  string     `gorm:"not null" json:"chat_gpt_model"`
	ChatGptToken  string     `gorm:"not null" json:"chat_gpt_token,omitempty"`
	ActivateAt    *time.Time `json:"-"`
	ActivationKey string     `json:"-"`
	CreatedAt     time.Time  `json:"-" sql:"DEFAULT:current_timestamp"`
}
