package repository

import (
	"fmt"
	"time"

	"gorm.io/gorm"

	"github.com/whites11/podcast-sync-server/internal/models"
)

type SubscriptionsRepository struct {
	db *gorm.DB
}

func NewSubscriptionsRepository(db *gorm.DB) (*SubscriptionsRepository, error) {
	err := db.AutoMigrate(&models.Subscription{})
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	return &SubscriptionsRepository{
		db: db,
	}, nil
}

func (d *SubscriptionsRepository) GetDeviceSubscriptionsSince(device models.Device, since time.Time) ([]models.Subscription, error) {
	var subscriptions []models.Subscription
	d.db.Unscoped().Where("device_id = ? AND (created_at > ? OR updated_at > ? OR deleted_At > ?)", device.ID, since, since, since).Find(&subscriptions)

	return subscriptions, nil
}

func (d *SubscriptionsRepository) CreateOrBumpSubscription(device models.Device, url string) error {
	// Bump timestamp to any row for this url/device ?
	update := d.db.Model(&models.Subscription{}).Where("device_id = ? AND url = ?", device.ID, url).Update("updated_at", time.Now())
	if update.Error != nil {
		return fmt.Errorf("%w", update.Error)
	}

	if update.RowsAffected == 0 {
		// create new subscription if no row matched
		result := d.db.Create(&models.Subscription{
			URL:    url,
			Device: device,
		})
		if result.Error != nil {
			return fmt.Errorf("%w", result.Error)
		}
	}

	return nil
}

func (d *SubscriptionsRepository) DeleteSubscription(device models.Device, url string) error {
	result := d.db.Where("device_id = ? AND url = ?", device.ID, url).Delete(&models.Subscription{})
	if result.Error != nil {
		return fmt.Errorf("%w", result.Error)
	}

	return nil
}

func (d *SubscriptionsRepository) GetMaxTimestamp(device models.Device) (time.Time, error) {
	maxTS := time.Unix(0, 0)

	subs, err := d.GetDeviceSubscriptionsSince(device, time.Unix(0, 0))
	if err != nil {
		return maxTS, err
	}

	for _, sub := range subs {
		if sub.UpdatedAt.After(maxTS) {
			maxTS = sub.UpdatedAt
		}
		if sub.DeletedAt.Valid && sub.DeletedAt.Time.After(maxTS) {
			maxTS = sub.DeletedAt.Time
		}
	}

	return maxTS, nil
}
