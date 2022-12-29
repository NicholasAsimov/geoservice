package store

import (
	"context"
	"fmt"
	"net/netip"
	"strings"

	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5"

	"github.com/nicholasasimov/geoservice/model"
)

type Store struct {
	DB *pgx.Conn
}

func New(db *pgx.Conn) *Store {
	return &Store{DB: db}
}

func (s *Store) UpsertRecords(ctx context.Context, records []model.GeoRecord) error {
	columns := []string{"ip_address", "country_code", "country", "city", "latitude", "longitude", "mystery_value"}

	tx, err := s.DB.Begin(ctx)
	if err != nil {
		return fmt.Errorf("can't start transaction: %w", err)
	}

	_, err = tx.Exec(ctx, `CREATE TEMPORARY TABLE _temp_upsert_georecords (LIKE georecords INCLUDING ALL) ON COMMIT DROP`)
	if err != nil {
		return fmt.Errorf("can't create temp table: %w", err)
	}

	copyCount, err := tx.CopyFrom(
		ctx,
		pgx.Identifier{"_temp_upsert_georecords"},
		columns,
		CopyFromRecords(records),
	)
	if err != nil {
		return fmt.Errorf("can't copy to temp table: %w", err)
	}

	_, err = tx.Exec(ctx, `INSERT INTO georecords SELECT * FROM _temp_upsert_georecords ON CONFLICT (ip_address) DO UPDATE SET `+buildSetSQL(columns))
	if err != nil {
		return fmt.Errorf("can't copy from temp table to actual table: %w", err)
	}

	err = tx.Commit(ctx)
	if err != nil {
		return fmt.Errorf("can't commit transaction: %w", err)
	}

	if copyCount != int64(len(records)) {
		return fmt.Errorf("incomplete db copy, records: %d, copied: %d", len(records), copyCount)
	}

	return nil
}

func buildSetSQL(columns []string) string {
	var out string
	for _, column := range columns {
		out += fmt.Sprintf("%s=EXCLUDED.%s, ", column, column)
	}
	return strings.TrimSuffix(out, ", ")
}

func (s *Store) GetRecord(ctx context.Context, addr netip.Addr) (model.GeoRecord, error) {
	var record model.GeoRecord

	err := pgxscan.Get(ctx, s.DB, &record, `SELECT * FROM georecords WHERE ip_address = $1`, addr.String())
	if err != nil {
		return record, fmt.Errorf("can't query db: %w", err)
	}
	return record, nil
}
