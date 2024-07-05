package models

import (
	"time"

	"gorm.io/gorm"
)

type EpisodeAction struct {
	ID       uint   `gorm:"primarykey"`
	Episode  string `gorm:"<-:create;uniqueIndex:eaud"`
	Action   string `gorm:"<-:create;uniqueIndex:eaud"`
	Started  uint
	Position uint
	Total    uint
	UserID   *uint `gorm:"<-:create;uniqueIndex:eaud"`
	DeviceID *uint `gorm:"<-:create;uniqueIndex:eaud"`

	User   User   `json:"-" gorm:"foreignKey:user_id;references:id"`
	Device Device `json:"-" gorm:"foreignKey:device_id;references:id"`

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}
