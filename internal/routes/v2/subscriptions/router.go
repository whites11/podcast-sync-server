package subscriptions

import (
	"github.com/go-chi/chi/v5"

	"github.com/whites11/podcast-sync-server/internal/dependencies"
	"github.com/whites11/podcast-sync-server/internal/middlewares"
)

type Router struct {
	deps *dependencies.Dependencies
}

func NewSubscriptionsRouter(deps *dependencies.Dependencies) *Router {
	return &Router{
		deps: deps,
	}
}

func (d *Router) Setup(r chi.Router) {
	r.With(middlewares.CheckAuthz(d.deps.DevicesRepository())).With(middlewares.BindAndValidate[getSubscriptionChangesRequest]).Get("/subscriptions/{username}/{deviceid}.json", d.getSubscriptionChangesHandler)
	r.With(middlewares.CheckAuthz(d.deps.DevicesRepository())).With(middlewares.BindAndValidate[uploadSubscriptionChangesRequest]).Post("/subscriptions/{username}/{deviceid}.json", d.uploadSubscriptionChangesHandler)
}
