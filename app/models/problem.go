package models

import "time"

type Problem struct {
	ID        uint      `gorm:"primarykey" json:"id"`
	Question  string    `gorm:"not null" json:"question"`
	CreatedAt time.Time `sql:"DEFAULT:current_timestamp" json:"-"`
	ThemeID   uint      `gorm:"not null" json:"theme_id" json:"-"`
	Theme     *Theme    `json:"theme,omitempty" json:"-"`
	Languages []string  `gorm:"not null" json:"languages"`
	IsTheory  bool      `gorm:"not null" json:"is_theory"`
}
