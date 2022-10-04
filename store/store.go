package store

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"

	"github.com/nicholasasimov/findhotel-assignment/model"
)

type Store struct {
	DB *pgx.Conn
}

func New(db *pgx.Conn) *Store {
	return &Store{DB: db}
}

func (s *Store) UpsertRecords(ctx context.Context, records []model.GeoRecord) error {
	copyCount, err := s.DB.CopyFrom(
		ctx,
		pgx.Identifier{"geo_records"},
		[]string{"ip_address", "country_code", "country", "city", "latitude", "longitude", "mystery_value"},
		CopyFromRecords(records),
	)
	if err != nil {
		return err
	}

	if copyCount != int64(len(records)) {
		return fmt.Errorf("incomplete db copy, records: %d, copied: %d", len(records), copyCount)
	}

	return nil
}
