package models

import "time"

type Problem struct {
	ID        uint      `gorm:"primarykey"`
	Question  string    `gorm:"not null" json:"question"`
	CreatedAt time.Time `sql:"DEFAULT:current_timestamp"`
	ThemeID   uint      `gorm:"not null" json:"theme_id"`
	Theme     *Theme    `json:"theme,omitempty"`
}
