package store

import (
	"context"
	"fmt"
	"net/netip"

	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5"

	"github.com/nicholasasimov/findhotel-assignment/model"
)

type Store struct {
	DB *pgx.Conn
}

func New(db *pgx.Conn) *Store {
	return &Store{DB: db}
}

// FIXME copy to temp table then upsert https://github.com/jackc/pgx/issues/992
func (s *Store) UpsertRecords(ctx context.Context, records []model.GeoRecord) error {
	copyCount, err := s.DB.CopyFrom(
		ctx,
		pgx.Identifier{"georecords"},
		[]string{"ip_address", "country_code", "country", "city", "latitude", "longitude", "mystery_value"},
		CopyFromRecords(records),
	)
	if err != nil {
		return fmt.Errorf("can't copy to temp table: %w", err)
	}

	if copyCount != int64(len(records)) {
		return fmt.Errorf("incomplete db copy, records: %d, copied: %d", len(records), copyCount)
	}

	return nil
}

func (s *Store) GetRecord(ctx context.Context, addr netip.Addr) (model.GeoRecord, error) {
	var record model.GeoRecord

	err := pgxscan.Get(ctx, s.DB, &record, `SELECT * FROM georecords WHERE ip_address = $1`, addr.String())
	if err != nil {
		return record, fmt.Errorf("can't query db: %w", err)
	}
	return record, nil
}
