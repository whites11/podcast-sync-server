package devices

import (
	"net/http"

	"github.com/whites11/podcast-sync-server/internal/errorrenderer"
	"github.com/whites11/podcast-sync-server/internal/key"
	"github.com/whites11/podcast-sync-server/internal/models"
)

type updateDeviceDataRequest struct {
	DeviceID string `pathvalue:"deviceid"`
	Username string `pathvalue:"username"`

	Caption string `json:"caption"`
	Type    string `json:"type"`
}

func (d *Router) updateDeviceDataHandler(w http.ResponseWriter, r *http.Request) {
	currentUser := r.Context().Value(key.CONTEXT_VALUE_CURRENT_USERNAME_KEY).(models.User)

	req := r.Context().Value(key.CONTEXT_VALUE_REQUEST_KEY).(updateDeviceDataRequest)

	_, err := d.devicesRepository.CreateDevice(models.Device{
		Slug:    req.DeviceID,
		Caption: req.Caption,
		Type:    req.Type,
		User:    currentUser,
	})
	if err != nil {
		errorrenderer.Render(w, r, err, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
