package models

type UserTheme struct {
	ID               uint   `gorm:"primarykey"`
	ThemeID          uint   `gorm:"not null" json:"theme_id"`
	Theme            Theme  `gorm:"foreignKey:ThemeID" json:"theme,omitempty"`
	UserID           uint   `gorm:"not null"`
	User             User   `gorm:"foreignKey:UserID" json:"-"`
	Score            uint   `gorm:"not null" json:"score"`
	ResolvedProblems string `json:"-"`
}
