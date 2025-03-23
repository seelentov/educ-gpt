package models

import "time"

type Problem struct {
	ID        uint      `gorm:"primarykey" json:"id,omitempty"`
	Question  string    `gorm:"not null" json:"question"`
	CreatedAt time.Time `sql:"DEFAULT:current_timestamp" json:"-"`
	ThemeID   uint      `gorm:"not null" json:"-"`
	Theme     *Theme    `json:"-"`
	Languages string    `gorm:"not null" json:"languages"`
	IsTheory  bool      `gorm:"not null" json:"is_theory"`
}
