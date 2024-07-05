package episodeactions

import (
	"net/http"

	"github.com/whites11/podcast-sync-server/internal/errorrenderer"
	"github.com/whites11/podcast-sync-server/internal/key"
	"github.com/whites11/podcast-sync-server/internal/models"
)

type uploadEpisodeActionsRequest struct {
	Actions []struct {
		Podcast    string `json:"podcast"`
		Episode    string `json:"episode"`
		DeviceSlug string `json:"device"`
		Action     string `json:"action"`
		Timestamp  string `json:"timestamp"`
		Started    uint   `json:"started"`
		Position   uint   `json:"position"`
		Total      uint   `json:"total"`
	} `json:"actions" jsonbodyfield:"actions"`
}

func (d *Router) uploadEpisodesActionsHandler(w http.ResponseWriter, r *http.Request) {
	currentUser := r.Context().Value(key.CONTEXT_VALUE_CURRENT_USERNAME_KEY).(models.User)
	req := r.Context().Value(key.CONTEXT_VALUE_REQUEST_KEY).(uploadEpisodeActionsRequest)

	actions := make([]*models.EpisodeAction, 0)

	devicesCache := make(map[string]*models.Device)

	for _, action := range req.Actions {
		deviceID := uint(0)
		if action.DeviceSlug != "" {
			if device, found := devicesCache[action.DeviceSlug]; found {
				deviceID = device.ID
			} else {
				var err error
				device, err = d.deps.DevicesRepository().FindBySlug(action.DeviceSlug)
				if err == nil {
					devicesCache[device.Slug] = device
					deviceID = device.ID
				}
			}
		}

		actions = append(actions, &models.EpisodeAction{
			Episode:  action.Episode,
			Action:   action.Action,
			Started:  action.Started,
			Position: action.Position,
			Total:    action.Total,
			UserID:   &currentUser.ID,
			DeviceID: &deviceID,
		})
	}

	err := d.deps.EpisodesActionsRepository().CreateBatch(actions)
	if err != nil {
		errorrenderer.Render(w, r, err, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
