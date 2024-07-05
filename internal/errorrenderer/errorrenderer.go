package errorrenderer

import (
	"net/http"

	"github.com/go-chi/render"
)

type errResponse struct {
	Err            error `json:"-"` // low-level runtime error
	HTTPStatusCode int   `json:"-"` // http response status code

	StatusText string `json:"status"`          // user-level status message
	AppCode    int64  `json:"code,omitempty"`  // application-specific error code
	ErrorText  string `json:"error,omitempty"` // application-level error message, for debugging
}

func (e *errResponse) Render(w http.ResponseWriter, r *http.Request) error {
	render.Status(r, e.HTTPStatusCode)
	return nil
}

func Render(w http.ResponseWriter, r *http.Request, err error, code int) {
	render.Render(w, r, &errResponse{
		Err:            err,
		HTTPStatusCode: code,
		StatusText:     err.Error(),
		ErrorText:      err.Error(),
	})
}
