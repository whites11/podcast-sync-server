package devices

import (
	"net/http"

	"github.com/go-chi/render"

	"github.com/whites11/podcast-sync-server/internal/errorrenderer"
	"github.com/whites11/podcast-sync-server/internal/key"
	"github.com/whites11/podcast-sync-server/internal/models"
)

type deviceResponse struct {
	models.Device
	SubscriptionsCount int `json:"subscriptions"`
}

type listDevicesResponse []deviceResponse

func (k *listDevicesResponse) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func (d *Router) listDevicesHandler(w http.ResponseWriter, r *http.Request) {
	currentUser := r.Context().Value(key.CONTEXT_VALUE_CURRENT_USERNAME_KEY).(models.User)

	devices, err := d.devicesRepository.GetUserDevices(currentUser)
	if err != nil {
		errorrenderer.Render(w, r, err, http.StatusInternalServerError)
		return
	}

	resp := listDevicesResponse{}
	for _, dev := range devices {
		resp = append(resp, deviceResponse{
			Device:             dev,
			SubscriptionsCount: len(dev.Subscriptions),
		})
	}

	if err := render.Render(w, r, &resp); err != nil {
		errorrenderer.Render(w, r, err, http.StatusInternalServerError)
		return
	}
}
