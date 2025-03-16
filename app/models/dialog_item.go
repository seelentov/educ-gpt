package models

import "time"

type DialogItem struct {
	ID        uint      `gorm:"primarykey" json:"id,omitempty"`
	Text      string    `gorm:"not null" json:"text,omitempty"`
	IsUser    bool      `gorm:"not null" json:"is_user"`
	DialogID  uint      `gorm:"not null"  json:"-"`
	Dialog    Dialog    `gorm:"foreignkey:DialogID" json:"-"`
	CreatedAt time.Time `sql:"DEFAULT:current_timestamp" json:"created_at"`
}
