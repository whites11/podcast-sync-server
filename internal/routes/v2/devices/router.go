package devices

import (
	"github.com/go-chi/chi/v5"

	"github.com/whites11/podcast-sync-server/internal/dependencies"
	"github.com/whites11/podcast-sync-server/internal/middlewares"
	"github.com/whites11/podcast-sync-server/internal/repository"
)

type Router struct {
	devicesRepository *repository.DevicesRepository
}

func NewDevicesRouter(deps *dependencies.Dependencies) *Router {
	return &Router{
		devicesRepository: deps.DevicesRepository(),
	}
}

func (d *Router) Setup(r chi.Router) {
	r.With(middlewares.CheckAuthz(d.devicesRepository)).Get("/devices/{username}.json", d.listDevicesHandler)
	r.With(middlewares.CheckAuthz(d.devicesRepository)).With(middlewares.BindAndValidate[updateDeviceDataRequest]).Post("/devices/{username}/{deviceid}.json", d.updateDeviceDataHandler)
}
