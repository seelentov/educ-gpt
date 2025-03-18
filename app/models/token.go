package models

import "time"

type TokenType int

const (
	TypeResetPassword TokenType = iota
	TypeChangeEmail
)

type Token struct {
	ID        uint      `gorm:"primarykey"`
	CreatedAt time.Time `sql:"DEFAULT:current_timestamp"`
	Key       string    `gorm:"not null"`
	UserID    uint      `gorm:"not null"`
	User      User      `gorm:"foreignkey:UserID"`
	Type      TokenType `gorm:"not null"`
	Data      string
}
