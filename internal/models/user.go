package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID           uint   `gorm:"primarykey"`
	Username     string `gorm:"uniqueIndex;<-:create"`
	PasswordHash string `json:"-"`
	Password     string `gorm:"-"`

	Devices []Device

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}
