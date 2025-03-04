package models

type Problem struct {
	ID       uint   `gorm:"primarykey"`
	Question string `gorm:"not null" json:"question"`
}
