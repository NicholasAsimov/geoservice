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
		return fmt.Errorf("copy was incomplete, records: %d, copied: %d", len(records), copyCount)
	}

	return nil
}

func CopyFromRecords(rows []model.GeoRecord) pgx.CopyFromSource {
	return &copyFromRecords{rows: rows, idx: -1}
}

type copyFromRecords struct {
	rows []model.GeoRecord
	idx  int
}

func (ctr *copyFromRecords) Next() bool {
	ctr.idx++
	return ctr.idx < len(ctr.rows)
}

func (ctr *copyFromRecords) Values() ([]any, error) {
	return []interface{}{
			ctr.rows[ctr.idx].IPAddress.String(),
			ctr.rows[ctr.idx].CountryCode,
			ctr.rows[ctr.idx].Country,
			ctr.rows[ctr.idx].City,
			ctr.rows[ctr.idx].Latitude,
			ctr.rows[ctr.idx].Longitude,
			ctr.rows[ctr.idx].MysteryValue,
		},
		nil
}

func (ctr *copyFromRecords) Err() error {
	return nil
}
