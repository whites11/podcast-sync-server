package routes

import (
	"github.com/go-chi/chi/v5"

	"github.com/whites11/podcast-sync-server/internal/dependencies"
	"github.com/whites11/podcast-sync-server/internal/middlewares"
	"github.com/whites11/podcast-sync-server/internal/routes/v2/auth"
	"github.com/whites11/podcast-sync-server/internal/routes/v2/devices"
	"github.com/whites11/podcast-sync-server/internal/routes/v2/episodeactions"
	"github.com/whites11/podcast-sync-server/internal/routes/v2/subscriptions"
)

const realm = "podcast-sync-server"

type V2Router struct {
	r chi.Router

	deps *dependencies.Dependencies
}

func NewV2Router(r chi.Router, deps *dependencies.Dependencies) *V2Router {
	return &V2Router{
		r:    r,
		deps: deps,
	}
}

func (v *V2Router) Setup() {
	v.r.Route("/2", func(r chi.Router) {
		r.Use(middlewares.RequireBasicAuth(v.deps, realm))
		auth.NewAuthRouter(v.deps).Setup(r)
		devices.NewDevicesRouter(v.deps).Setup(r)
		subscriptions.NewSubscriptionsRouter(v.deps).Setup(r)
		episodeactions.NewEpisodesRouter(v.deps).Setup(r)
	})
}
