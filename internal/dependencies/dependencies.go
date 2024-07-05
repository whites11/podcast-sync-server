package dependencies

import (
	"github.com/gorilla/securecookie"

	"github.com/whites11/podcast-sync-server/internal/repository"
)

type Dependencies struct {
	devicesRepository        *repository.DevicesRepository
	episodeActionsRepository *repository.EpisodeActionsRepository
	subscriptionsRepository  *repository.SubscriptionsRepository
	usersRepository          *repository.UsersRepository
	secureCookie             *securecookie.SecureCookie
}

func New(devicesRepository *repository.DevicesRepository, episodeActionsRepository *repository.EpisodeActionsRepository, subscriptionsRepository *repository.SubscriptionsRepository, usersRepository *repository.UsersRepository, secureCookie *securecookie.SecureCookie) *Dependencies {
	return &Dependencies{
		devicesRepository:        devicesRepository,
		episodeActionsRepository: episodeActionsRepository,
		subscriptionsRepository:  subscriptionsRepository,
		usersRepository:          usersRepository,
		secureCookie:             secureCookie,
	}
}

func (d *Dependencies) DevicesRepository() *repository.DevicesRepository {
	return d.devicesRepository
}

func (d *Dependencies) EpisodesActionsRepository() *repository.EpisodeActionsRepository {
	return d.episodeActionsRepository
}

func (d *Dependencies) SubscriptionsRepository() *repository.SubscriptionsRepository {
	return d.subscriptionsRepository
}

func (d *Dependencies) UsersRepository() *repository.UsersRepository {
	return d.usersRepository
}

func (d *Dependencies) SecureCookie() *securecookie.SecureCookie { return d.secureCookie }
