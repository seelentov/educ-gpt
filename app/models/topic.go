package models

type Topic struct {
	ID     uint     `gorm:"primarykey" json:"id"`
	Title  string   `gorm:"unique, not null" json:"title"`
	Themes []*Theme `gorm:"foreignKey:TopicID" json:"themes,omitempty"`
	Score  uint     `gorm:"-" json:"scores,omitempty"`
}
