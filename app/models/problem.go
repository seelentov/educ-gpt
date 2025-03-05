package models

type Problem struct {
	ID       uint   `gorm:"primarykey"`
	Question string `gorm:"not null" json:"question"`
	ThemeID  uint   `gorm:"not null" json:"theme_id"`
	Theme    *Theme `json:"theme,omitempty"`
}
