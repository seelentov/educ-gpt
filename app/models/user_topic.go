package models

type UserTopic struct {
	ID      uint `gorm:"primary_key"`
	UserID  uint `gorm:"not null"`
	User    User
	TopicID uint `gorm:"not null"`
	Topic   Topic
}
