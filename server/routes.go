package server

import (
	"github.com/go-chi/chi/v5"
)

func (s *Server) Routes() chi.Router {
	r := chi.NewRouter()

	r.Get("/geolocation", s.GetGeolocation)

	return r
}
