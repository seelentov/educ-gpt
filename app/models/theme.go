package models

type Theme struct {
	ID               uint        `gorm:"primarykey" json:"id"`
	Title            string      `gorm:"not null" json:"title"`
	TopicID          uint        `gorm:"not null" json:"topic_id"`
	Topic            *Topic      `json:"topic,omitempty"`
	UserThemes       []UserTheme `gorm:"foreignKey:ThemeID" json:"-"` // Добавлено
	Score            uint        `gorm:"-" json:"scores,omitempty"`
	ResolvedProblems string      `gorm:"-" json:"-"`
}
