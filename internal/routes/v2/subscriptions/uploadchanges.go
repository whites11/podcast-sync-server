package subscriptions

import (
	"net/http"

	"github.com/go-chi/render"

	"github.com/whites11/podcast-sync-server/internal/errorrenderer"
	"github.com/whites11/podcast-sync-server/internal/key"
	"github.com/whites11/podcast-sync-server/internal/models"
)

type uploadSubscriptionChangesRequest struct {
	Add    []string `json:"add"`
	Remove []string `json:"remove"`
}

type uploadSubscriptionChangesResponse struct {
	Timestamp  int64      `json:"timestamp"`
	UpdateURLs [][]string `json:"update_urls"`
}

func (k *uploadSubscriptionChangesResponse) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func (d *Router) uploadSubscriptionChangesHandler(w http.ResponseWriter, r *http.Request) {
	currentDevice := r.Context().Value(key.CONTEXT_VALUE_CURRENT_DEVICE_KEY).(models.Device)
	req := r.Context().Value(key.CONTEXT_VALUE_REQUEST_KEY).(uploadSubscriptionChangesRequest)

	for _, added := range req.Add {
		err := d.deps.SubscriptionsRepository().CreateOrBumpSubscription(currentDevice, added)
		if err != nil {
			errorrenderer.Render(w, r, err, http.StatusInternalServerError)
			return
		}
	}

	for _, deleted := range req.Remove {
		err := d.deps.SubscriptionsRepository().DeleteSubscription(currentDevice, deleted)
		if err != nil {
			errorrenderer.Render(w, r, err, http.StatusInternalServerError)
			return
		}
	}

	ts, err := d.deps.SubscriptionsRepository().GetMaxTimestamp(currentDevice)
	if err != nil {
		errorrenderer.Render(w, r, err, http.StatusInternalServerError)
		return
	}

	resp := &uploadSubscriptionChangesResponse{
		Timestamp:  ts.Unix(),
		UpdateURLs: make([][]string, 0),
	}

	if err := render.Render(w, r, resp); err != nil {
		errorrenderer.Render(w, r, err, http.StatusInternalServerError)
		return
	}
}
