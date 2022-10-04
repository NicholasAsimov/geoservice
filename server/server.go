package server

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/netip"

	"github.com/jackc/pgx/v5"
	"github.com/rs/zerolog"

	"github.com/nicholasasimov/findhotel-assignment/store"
)

type Server struct {
	Log   zerolog.Logger
	DB    *pgx.Conn
	Store *store.Store
}

func New(log zerolog.Logger, db *pgx.Conn, s *store.Store) *Server {
	return &Server{
		Log:   log,
		DB:    db,
		Store: s,
	}
}

type APIError struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
}

func (s *Server) apiError(w http.ResponseWriter, err string, statusCode int) {
	e := APIError{
		Message: err,
		Code:    statusCode,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	if err := json.NewEncoder(w).Encode(e); err != nil {
		s.Log.Error().Err(err).Msg("can't write response")
	}
}

func (s *Server) GetGeolocation(w http.ResponseWriter, r *http.Request) {
	addr, err := netip.ParseAddr(r.URL.Query().Get("ip"))
	if err != nil {
		s.Log.Warn().Err(err).Msg("can't parse ip")
		s.apiError(w, "can't parse ip", http.StatusBadRequest)
		return
	}

	record, err := s.Store.GetRecord(r.Context(), addr)
	if errors.Is(err, pgx.ErrNoRows) {
		s.apiError(w, "no geolocation data for this ip", http.StatusNotFound)
		return
	}
	if err != nil {
		s.Log.Error().Err(err).Msg("can't get geolocation data")
		s.apiError(w, "can't get geolocation data", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)

	if err := json.NewEncoder(w).Encode(record); err != nil {
		s.Log.Error().Err(err).Msg("can't write response")
	}
}
