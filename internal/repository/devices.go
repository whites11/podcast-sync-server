package repository

import (
	"fmt"

	"gorm.io/gorm"

	"github.com/whites11/podcast-sync-server/internal/models"
)

type DevicesRepository struct {
	db *gorm.DB
}

func NewDevicesRepository(db *gorm.DB) (*DevicesRepository, error) {
	err := db.AutoMigrate(&models.Device{})
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	return &DevicesRepository{
		db: db,
	}, nil
}

func (d *DevicesRepository) GetUserDevices(user models.User) ([]models.Device, error) {
	var devices []models.Device
	result := d.db.Where("user_id = ?", user.ID).Preload("Subscriptions").Find(&devices)
	if result.Error != nil {
		return nil, fmt.Errorf("error getting user devices: %w", result.Error)
	}

	return devices, nil
}

func (d *DevicesRepository) FindBySlug(slug string) (*models.Device, error) {
	var device models.Device
	result := d.db.First(&device, "slug = ?", slug)
	if result.Error != nil {
		return nil, fmt.Errorf("device not found")
	}

	return &device, nil
}

func (d *DevicesRepository) CreateDevice(device models.Device) (*models.Device, error) {
	result := d.db.Create(&device)

	if result.Error != nil {
		return nil, fmt.Errorf("%w", result.Error)
	}

	return &device, nil
}
