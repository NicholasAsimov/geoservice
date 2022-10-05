package csvparse

import (
	"encoding/csv"
	"fmt"
	"io"
	"net/netip"

	"github.com/jszwec/csvutil"
	"golang.org/x/exp/maps"

	"github.com/nicholasasimov/findhotel-assignment/model"
)

// ParseCSV parses the CSV content from the io.Reader until EOF, returning
// successfully parsed records and the number of skipped records.
// It validates each record in a streaming manner using validateFunc.
// Duplicated records are discarded, first record for a given IP is used.
func ParseCSV(r io.Reader, validateFunc func(model.GeoRecord) bool) ([]model.GeoRecord, int, error) {
	dec, err := csvutil.NewDecoder(csvReader(r))
	if err != nil {
		return nil, 0, fmt.Errorf("can't create csv decoder: %w", err)
	}

	records := make(map[netip.Addr]model.GeoRecord)
	var skipped int
	for {
		var record model.GeoRecord

		err := dec.Decode(&record)
		if err == io.EOF {
			break
		}

		_, exists := records[record.IPAddress]
		if err != nil || exists || !validateFunc(record) {
			skipped += 1
			continue
		}

		records[record.IPAddress] = record
	}

	return maps.Values(records), skipped, nil
}

func csvReader(in io.Reader) *csv.Reader {
	r := csv.NewReader(in)
	r.LazyQuotes = true
	r.TrimLeadingSpace = true
	return r
}
