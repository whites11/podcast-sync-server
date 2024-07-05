package auth

import (
	"fmt"
	"net/http"

	"github.com/whites11/podcast-sync-server/internal/errorrenderer"
	"github.com/whites11/podcast-sync-server/internal/key"
)

type loginRequest struct {
	Username string `pathvalue:"username"`
}

func (a *Router) loginHandler(w http.ResponseWriter, r *http.Request) {
	req := r.Context().Value(key.CONTEXT_VALUE_REQUEST_KEY).(loginRequest)

	value := map[string]string{
		key.COOKIE_VALUE_KEY_USERNAME: req.Username,
	}

	encoded, err := a.deps.SecureCookie().Encode(key.COOKIE_NAME, value)
	if err != nil {
		fmt.Sprintf("error encoding cookie: %s", err)
		errorrenderer.Render(w, r, err, http.StatusInternalServerError)
		return
	} else {
		cookie := http.Cookie{
			Name:     key.COOKIE_NAME,
			Value:    encoded,
			Path:     "/",
			MaxAge:   3600,
			HttpOnly: true,
			Secure:   false,
			SameSite: http.SameSiteNoneMode,
		}

		http.SetCookie(w, &cookie)
	}

	w.WriteHeader(http.StatusOK)
}
