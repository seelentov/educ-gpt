package models

import "time"

type Theme struct {
	ID               uint        `gorm:"primarykey" json:"id"`
	Title            string      `gorm:"not null" json:"title"`
	TopicID          uint        `gorm:"not null" json:"topic_id"`
	Topic            *Topic      `json:"topic,omitempty"`
	UserThemes       []UserTheme `gorm:"foreignKey:ThemeID;constraint:OnDelete:CASCADE;" json:"-"`
	Score            uint        `gorm:"-" json:"scores,omitempty"`
	ResolvedProblems string      `gorm:"-" json:"-"`
	CreatedAt        time.Time   `json:"-" sql:"DEFAULT:current_timestamp"`
}
