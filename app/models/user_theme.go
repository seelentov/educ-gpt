package models

type UserTheme struct {
	ID               uint   `gorm:"primarykey"`
	Theme            *Theme `json:"theme,omitempty"`
	ThemeID          uint   `gorm:"not null" json:"theme_id"`
	User             *Topic `json:"user,omitempty"`
	UserID           uint   `gorm:"not null" json:"user_id"`
	Score            uint   `gorm:"not null" json:"score"`
	ResolvedProblems string `gorm:"-" json:"-"`
}
