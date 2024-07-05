package middlewares

import (
	"context"
	"fmt"
	"net/http"

	"github.com/whites11/podcast-sync-server/internal/dependencies"
	"github.com/whites11/podcast-sync-server/internal/errorrenderer"
	"github.com/whites11/podcast-sync-server/internal/key"
	"github.com/whites11/podcast-sync-server/internal/models"
)

// RequireBasicAuth implements a simple middleware handler for adding basic http auth to a route.
func RequireBasicAuth(deps *dependencies.Dependencies, realm string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			var user *models.User

			// Check cookie
			c, err := r.Cookie(key.COOKIE_NAME)
			if err == nil && c.Valid() == nil {
				// Cookie found
				value := make(map[string]string)
				err = deps.SecureCookie().Decode(key.COOKIE_NAME, c.Value, &value)
				if err != nil {
					// Invalid cookie.
					errorrenderer.Render(w, r, fmt.Errorf("invalid cookie"), http.StatusUnauthorized)
					return
				}

				username := value[key.COOKIE_VALUE_KEY_USERNAME]
				if username != "" {
					user, err = deps.UsersRepository().FindByUsername(username)
					if err != nil {
						// Invalid cookie.
						errorrenderer.Render(w, r, fmt.Errorf("invalid cookie"), http.StatusUnauthorized)
						return
					}
				}
			}

			if user == nil {
				username, password, ok := r.BasicAuth()
				if !ok {
					basicAuthFailed(w, realm)
					return
				}

				user, err = deps.UsersRepository().FindByCredentials(username, password)
				if err != nil {
					basicAuthFailed(w, realm)
					return
				}
			}

			ctxWithUser := context.WithValue(r.Context(), key.CONTEXT_VALUE_CURRENT_USERNAME_KEY, *user)

			next.ServeHTTP(w, r.WithContext(ctxWithUser))
		})
	}
}

func basicAuthFailed(w http.ResponseWriter, realm string) {
	w.Header().Add("WWW-Authenticate", fmt.Sprintf(`Basic realm="%s"`, realm))
	w.WriteHeader(http.StatusUnauthorized)
}
