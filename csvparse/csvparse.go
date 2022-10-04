package csvparse

import (
	"encoding/csv"
	"fmt"
	"io"

	"github.com/jszwec/csvutil"

	"github.com/nicholasasimov/findhotel-assignment/model"
)

// ParseCSV parses the CSV content from the io.Reader until EOF, returning
// successfully parsed records and the number of skipped records.
// It validates each record in a streaming manner using validateFunc.
func ParseCSV(r io.Reader, validateFunc func(model.Record) bool) ([]model.Record, int, error) {
	dec, err := csvutil.NewDecoder(csvReader(r))
	if err != nil {
		return nil, 0, fmt.Errorf("can't create csv decoder: %w", err)
	}

	var records []model.Record
	var skipped int
	for {
		var record model.Record

		err := dec.Decode(&record)
		if err == io.EOF {
			break
		}

		if err != nil || !validateFunc(record) {
			skipped += 1
			continue
		}

		records = append(records, record)
	}

	return records, skipped, nil
}

func csvReader(in io.Reader) *csv.Reader {
	r := csv.NewReader(in)
	r.LazyQuotes = true
	return r
}
