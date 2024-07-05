package subscriptions

import (
	"net/http"
	"time"

	"github.com/go-chi/render"

	"github.com/whites11/podcast-sync-server/internal/errorrenderer"
	"github.com/whites11/podcast-sync-server/internal/key"
	"github.com/whites11/podcast-sync-server/internal/models"
)

type getSubscriptionChangesRequest struct {
	Since int64 `querystring:"since" validate:"gte=0"`
}

type getSubscriptionChangesResponse struct {
	Add       []string `json:"add"`
	Remove    []string `json:"remove"`
	Timestamp uint     `json:"timestamp"`
}

func (k *getSubscriptionChangesResponse) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func (d *Router) getSubscriptionChangesHandler(w http.ResponseWriter, r *http.Request) {
	currentDevice := r.Context().Value(key.CONTEXT_VALUE_CURRENT_DEVICE_KEY).(models.Device)
	req := r.Context().Value(key.CONTEXT_VALUE_REQUEST_KEY).(getSubscriptionChangesRequest)

	changes, err := d.deps.SubscriptionsRepository().GetDeviceSubscriptionsSince(currentDevice, time.Unix(req.Since, 0))
	if err != nil {
		errorrenderer.Render(w, r, err, http.StatusInternalServerError)
		return
	}

	add := []string{}
	remove := []string{}

	var maxTS int64
	for _, sub := range changes {
		if sub.CreatedAt.Unix() > maxTS {
			maxTS = sub.CreatedAt.Unix()
		}
		if sub.UpdatedAt.Unix() > maxTS {
			maxTS = sub.UpdatedAt.Unix()
		}
		if sub.DeletedAt.Valid {
			if sub.DeletedAt.Time.Unix() > maxTS {
				maxTS = sub.DeletedAt.Time.Unix()
			}
			remove = append(remove, sub.URL)
		} else {
			add = append(add, sub.URL)
		}
	}

	resp := &getSubscriptionChangesResponse{
		Add:       add,
		Remove:    remove,
		Timestamp: uint(maxTS),
	}

	if err := render.Render(w, r, resp); err != nil {
		errorrenderer.Render(w, r, err, http.StatusInternalServerError)
		return
	}
}
