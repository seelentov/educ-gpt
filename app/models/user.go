package models

import "time"

type User struct {
	ID            uint         `gorm:"primarykey" json:"id"`
	Name          string       `gorm:"unique;not null" json:"name"`
	Email         string       `gorm:"unique;not null" json:"email"`
	Password      string       `gorm:"not null" json:"-"`
	AvatarUrl     string       `json:"avatar_url"`
	ChatGptModel  string       `gorm:"not null" json:"chat_gpt_model"`
	ChatGptToken  string       `gorm:"not null" json:"chat_gpt_token,omitempty"`
	ActivationKey string       `json:"-"`
	ActivateAt    *time.Time   `json:"-"`
	CreatedAt     time.Time    `json:"-" sql:"DEFAULT:current_timestamp"`
	Roles         []*Role      `gorm:"many2many:user_roles;constraint:OnDelete:CASCADE" json:"-"`
	Themes        []*UserTheme `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE;" json:"-"` // Уточнено
	Dialogs       []*Dialog    `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE;" json:"-"`
}
