package server

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/netip"

	"github.com/jackc/pgx/v5"
	"github.com/rs/zerolog"

	"github.com/nicholasasimov/geoservice/config"
	"github.com/nicholasasimov/geoservice/csvparse"
	"github.com/nicholasasimov/geoservice/model"
	"github.com/nicholasasimov/geoservice/store"
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

func (s *Server) IngestData(w http.ResponseWriter, r *http.Request) {
	maxUploadSize := 1024 * 1024 * config.Config.Server.MaxUploadMB

	r.Body = http.MaxBytesReader(w, r.Body, maxUploadSize)
	if err := r.ParseMultipartForm(maxUploadSize); err != nil {
		s.apiError(w, "file too big", http.StatusBadRequest)
		return
	}

	file, fileHeader, err := r.FormFile("file")
	if err != nil {
		s.apiError(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer file.Close()

	validate := func(r model.GeoRecord) bool {
		return r.IPAddress.IsValid() && r.City != "" && r.Country != "" && r.CountryCode != ""
	}

	s.Log.Info().Str("file", fileHeader.Filename).Msg("parsing file")

	records, skipped, err := csvparse.ParseCSV(file, validate)
	if err != nil {
		s.Log.Error().Err(err).Msg("can't parse csv file")
		s.apiError(w, "can't parse csv file", http.StatusInternalServerError)
		return
	}

	s.Log.Info().Int("records", len(records)).Msg("persisting to db")

	err = s.Store.UpsertRecords(r.Context(), records)
	if err != nil {
		s.Log.Error().Err(err).Msg("can't save records in db")
		s.apiError(w, "can't save records in db", http.StatusInternalServerError)
		return
	}

	type resp struct {
		Records  int `json:"records"`
		Accepted int `json:"accepted"`
		Skipped  int `json:"skipped"`
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)

	if err := json.NewEncoder(w).Encode(resp{
		Records:  len(records) + skipped,
		Accepted: len(records),
		Skipped:  skipped,
	}); err != nil {
		s.Log.Error().Err(err).Msg("can't write response")
	}
}
