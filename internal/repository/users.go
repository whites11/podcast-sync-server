package repository

import (
	"errors"
	"fmt"
	"log"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"github.com/whites11/podcast-sync-server/internal/models"
)

type UsersRepository struct {
	db *gorm.DB
}

func NewUsersRepository(db *gorm.DB) (*UsersRepository, error) {
	err := db.AutoMigrate(&models.User{})
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	return &UsersRepository{
		db: db,
	}, nil
}

func (u *UsersRepository) FindByCredentials(username, pass string) (*models.User, error) {
	user, err := u.FindByUsername(username)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("user not found")
		}

		return nil, fmt.Errorf("unexpected error %w", err)
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(pass))
	if err != nil {
		return nil, fmt.Errorf("invalid password")
	}

	return user, nil
}

func (u *UsersRepository) FindByUsername(username string) (*models.User, error) {
	var user models.User

	err := u.db.First(&user, "username = ?", username).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("user not found")
	}

	return &user, nil
}

func (u *UsersRepository) CreateUser(user models.User) (*models.User, error) {
	if user.Password == "" {
		return nil, fmt.Errorf("cannot create user without a password")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.MinCost)
	if err != nil {
		log.Println(err)
	}

	user.PasswordHash = string(hash)

	result := u.db.Create(&user)
	if result.Error != nil {
		return nil, fmt.Errorf("error creating user: %w", result.Error)
	}

	return &user, nil
}
