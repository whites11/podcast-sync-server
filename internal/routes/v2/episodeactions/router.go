package episodeactions

import (
	"github.com/go-chi/chi/v5"

	"github.com/whites11/podcast-sync-server/internal/dependencies"
	"github.com/whites11/podcast-sync-server/internal/middlewares"
)

type Router struct {
	deps *dependencies.Dependencies
}

func NewEpisodesRouter(deps *dependencies.Dependencies) *Router {
	return &Router{
		deps: deps,
	}
}

func (d *Router) Setup(r chi.Router) {
	r.With(middlewares.CheckAuthz(d.deps.DevicesRepository())).With(middlewares.BindAndValidate[getEpisodeActionsRequest]).Get("/episodes/{username}.json", d.getEpisodesActionsHandler)
	r.With(middlewares.CheckAuthz(d.deps.DevicesRepository())).With(middlewares.BindAndValidate[uploadEpisodeActionsRequest]).Post("/episodes/{username}.json", d.uploadEpisodesActionsHandler)
}
