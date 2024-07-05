package models

import (
	"time"

	"gorm.io/gorm"
)

type Subscription struct {
	ID       uint   `gorm:"primarykey"`
	URL      string `json:"url" gorm:"<-:create"`
	UserID   uint   `json:"-" gorm:"<-:create"`
	DeviceID uint   `json:"-" gorm:"<-:create"`

	User   User   `json:"-" gorm:"foreignKey:user_id;references:id"`
	Device Device `json:"-" gorm:"foreignKey:device_id;references:id"`

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}
