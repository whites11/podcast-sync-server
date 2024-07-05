package models

import (
	"time"

	"gorm.io/gorm"
)

type Device struct {
	ID      uint   `json:"-" gorm:"primarykey"`
	Slug    string `json:"id" gorm:"uniqueIndex:deviceslug;index;<-:create"`
	Caption string `json:"caption"`
	Type    string `json:"type;<-:create"`
	UserID  uint   `json:"-" gorm:"uniqueIndex:deviceslug;<-:create"`

	User          User           `json:"-" gorm:"foreignKey:user_id;references:id"`
	Subscriptions []Subscription `json:"-"`

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}
