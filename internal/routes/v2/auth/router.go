package auth

import (
	"github.com/go-chi/chi/v5"

	"github.com/whites11/podcast-sync-server/internal/dependencies"
	"github.com/whites11/podcast-sync-server/internal/middlewares"
)

type Router struct {
	deps *dependencies.Dependencies
}

func NewAuthRouter(deps *dependencies.Dependencies) *Router {
	return &Router{
		deps: deps,
	}
}

func (a *Router) Setup(r chi.Router) {
	r.Route("/auth/{username}", func(r chi.Router) {
		r.With(middlewares.CheckAuthz(a.deps.DevicesRepository())).With(middlewares.BindAndValidate[loginRequest]).Post("/login.json", a.loginHandler)
		r.With(middlewares.CheckAuthz(a.deps.DevicesRepository())).Post("/logout.json", a.logoutHandler)
	})
}
