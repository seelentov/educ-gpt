package models

import "time"

type ResetToken struct {
	ID        uint      `gorm:"primarykey"`
	CreatedAt time.Time `sql:"DEFAULT:current_timestamp"`
	Key       string    `gorm:"not null"`
	UserID    uint      `gorm:"not null"`
	User      User      `gorm:"foreignkey:UserID"`
}
