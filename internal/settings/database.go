package settings

import (
	"errors"
	"fmt"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Setting struct {
	Key   string `gorm:"uniqueIndex"`
	Value string
}

type Database struct {
	db *gorm.DB
}

func NewDatabaseStorage(db *gorm.DB) (*Database, error) {
	err := db.AutoMigrate(&Setting{})
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	return &Database{
		db: db,
	}, nil
}

func (d *Database) Get(name string) (string, error) {
	setting := Setting{}
	tx := d.db.Where("key = ?", name).First(&setting)
	if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
		return "", nil
	} else if tx.Error != nil {
		return "", fmt.Errorf("%w", tx.Error)
	}

	return setting.Value, nil
}

func (d *Database) Set(name string, value string) error {
	d.db.Clauses(clause.OnConflict{
		UpdateAll: true,
	}).Create(&Setting{
		Key:   name,
		Value: value,
	})
	return nil
}
