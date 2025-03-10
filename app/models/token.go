package models

import "time"

type Type string

const (
	TypeResetPassword Type = "reset_password"
	TypeChangeEmail   Type = "change_email"
)

type Token struct {
	ID        uint      `gorm:"primarykey"`
	CreatedAt time.Time `sql:"DEFAULT:current_timestamp"`
	Key       string    `gorm:"not null"`
	UserID    uint      `gorm:"not null"`
	User      User      `gorm:"foreignkey:UserID"`
	Type      Type      `gorm:"not null"`
	Data      string
}
