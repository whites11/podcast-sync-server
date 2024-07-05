package episodeactions

import (
	"net/http"
	"time"

	"github.com/go-chi/render"

	"github.com/whites11/podcast-sync-server/internal/errorrenderer"
	"github.com/whites11/podcast-sync-server/internal/key"
	"github.com/whites11/podcast-sync-server/internal/models"
)

type getEpisodeActionsRequest struct {
	Podcast    string `querystring:"podcast"`
	Device     string `querystring:"device"`
	Since      int64  `querystring:"since" validate:"gte=0"`
	Aggregated bool   `querystring:"aggregated"`
}

type getEpisodesActionsResponse struct {
	Actions   []models.EpisodeAction `json:"actions"`
	Timestamp int64                  `json:"timestamp"`
}

func (k *getEpisodesActionsResponse) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func (d *Router) getEpisodesActionsHandler(w http.ResponseWriter, r *http.Request) {
	currentUser := r.Context().Value(key.CONTEXT_VALUE_CURRENT_USERNAME_KEY).(models.User)

	req := r.Context().Value(key.CONTEXT_VALUE_REQUEST_KEY).(getEpisodeActionsRequest)

	actions, err := d.deps.EpisodesActionsRepository().GetEpisodeActions(currentUser, time.Unix(req.Since, 0))
	if err != nil {
		errorrenderer.Render(w, r, err, http.StatusInternalServerError)
		return
	}

	resp := &getEpisodesActionsResponse{
		Actions:   actions,
		Timestamp: 0,
	}

	if err := render.Render(w, r, resp); err != nil {
		errorrenderer.Render(w, r, err, http.StatusInternalServerError)
		return
	}
}
