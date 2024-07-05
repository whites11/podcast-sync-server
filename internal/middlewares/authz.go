package middlewares

import (
	"context"
	"fmt"
	"net/http"

	"github.com/whites11/podcast-sync-server/internal/errorrenderer"
	"github.com/whites11/podcast-sync-server/internal/key"
	"github.com/whites11/podcast-sync-server/internal/models"
	"github.com/whites11/podcast-sync-server/internal/repository"
)

func CheckAuthz(devicesRepository *repository.DevicesRepository) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			var err error
			currentUser := r.Context().Value(key.CONTEXT_VALUE_CURRENT_USERNAME_KEY).(models.User)

			if currentUser.Username != r.PathValue(key.PATH_VALUE_USERNAME_KEY) {
				errorrenderer.Render(w, r, fmt.Errorf("you are not authorized to perform this action"), http.StatusUnauthorized)
				return
			}

			ctx := r.Context()

			deviceId := r.PathValue(key.PATH_VALUE_DEVICE_ID_KEY)
			var currentDevice *models.Device
			if deviceId != "" {
				currentDevice, err = devicesRepository.FindBySlug(deviceId)
				if err != nil && err.Error() != "device not found" {
					errorrenderer.Render(w, r, err, http.StatusInternalServerError)
					return
				}
			}

			if currentDevice != nil {
				if currentDevice.UserID != currentUser.ID {
					errorrenderer.Render(w, r, fmt.Errorf("you are not authorized to perform this action"), http.StatusUnauthorized)
					return
				}

				ctx = context.WithValue(r.Context(), key.CONTEXT_VALUE_CURRENT_DEVICE_KEY, *currentDevice)
			}

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
