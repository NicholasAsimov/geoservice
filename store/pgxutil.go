package store

import (
	"github.com/jackc/pgx/v5"

	"github.com/nicholasasimov/findhotel-assignment/model"
)

// CopyFromRecords allows pgx to directly copy GeoRecord struct without
// additional allocations for transforming to [][]interface.
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
