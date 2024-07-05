package db

import (
	"fmt"
	"os"
	"path/filepath"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type SqliteFactory struct {
	path string
}

func NewSqliteFactory() (*SqliteFactory, error) {
	return &SqliteFactory{
		path: "data/db.sqlite",
	}, nil
}

func (f *SqliteFactory) WithFilePath(path string) *SqliteFactory {
	f.path = path
	return f
}

func (f *SqliteFactory) Build() (*gorm.DB, error) {
	if filepath.Dir(f.path) != "" {
		err := os.MkdirAll(filepath.Dir(f.path), 0755)
		if err != nil {
			return nil, fmt.Errorf("%w", err)
		}
	}

	db, err := gorm.Open(sqlite.Open(f.path), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	return db, nil
}
