package routes

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

type ProbesRouter struct {
}

func NewProbesRouter() *ProbesRouter {
	return &ProbesRouter{}
}

func (p *ProbesRouter) Setup(r chi.Router) {
	r.Get("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	r.Get("/readyz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
}
