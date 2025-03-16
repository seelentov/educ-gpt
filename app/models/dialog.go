package models

type Dialog struct {
	ID          uint          `gorm:"primarykey" json:"id,omitempty"`
	UserID      uint          `gorm:"not null" json:"-"`
	User        User          `gorm:"foreignkey:UserID" json:"-"`
	DialogItems []*DialogItem `gorm:"foreignKey:DialogID;constraint:OnDelete:CASCADE;" json:"dialog_items"`
}
