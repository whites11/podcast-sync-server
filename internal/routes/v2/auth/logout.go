package auth

import (
	"net/http"

	"github.com/whites11/podcast-sync-server/internal/key"
)

func (a *Router) logoutHandler(w http.ResponseWriter, r *http.Request) {
	cookie := http.Cookie{
		Name:   key.COOKIE_NAME,
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	}

	http.SetCookie(w, &cookie)

	w.WriteHeader(http.StatusOK)
}
