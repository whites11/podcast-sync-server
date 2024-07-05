package repository

import (
	"fmt"
	"time"

	"gorm.io/gorm"

	"github.com/whites11/podcast-sync-server/internal/models"
)

type EpisodeActionsRepository struct {
	db *gorm.DB
}

func NewEpisodeActionsRepository(db *gorm.DB) (*EpisodeActionsRepository, error) {
	err := db.AutoMigrate(&models.EpisodeAction{})
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	return &EpisodeActionsRepository{
		db: db,
	}, nil
}

func (e *EpisodeActionsRepository) GetEpisodeActions(user models.User, since time.Time) ([]models.EpisodeAction, error) {
	var actions []models.EpisodeAction
	result := e.db.Where("user_id = ? AND updated_at > ?", user.ID, since.Unix()).Find(&actions)
	if result.Error != nil {
		return nil, fmt.Errorf("%w", result.Error)
	}

	return actions, nil
}

func (e *EpisodeActionsRepository) CreateBatch(actions []*models.EpisodeAction) error {
	err := e.db.Transaction(func(tx *gorm.DB) error {
		for _, action := range actions {
			if err := tx.Create(action).Error; err != nil {
				return err
			}
		}

		return nil
	})
	if err != nil {
		return fmt.Errorf("%w", err)
	}

	return nil
}
