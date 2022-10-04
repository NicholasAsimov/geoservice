package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"log"
	"net/netip"
	"os"

	"github.com/davecgh/go-spew/spew"
	"github.com/jszwec/csvutil"
)

type Record struct {
	IPAddress    netip.Addr `csv:"ip_address"`
	CountryCode  string     `csv:"country_code"`
	Country      string     `csv:"country"`
	City         string     `csv:"city"`
	Latitude     float64    `csv:"latitude"`
	Longitude    float64    `csv:"longitude"`
	MysteryValue float64    `csv:"mystery_value"`
}

func main() {
	var path string

	flag.StringVar(&path, "path", "", "path to the bank statement")
	flag.Parse()

	file, err := os.Open(path)
	if err != nil {
		log.Printf("can't open csv file: %s", err)
		return
	}
	defer file.Close()

	validate := func(r Record) bool {
		return r.IPAddress.IsValid() && r.City != "" && r.Country != "" && r.CountryCode != ""
	}

	records, skipped, err := ParseCSV(file, validate)
	if err != nil {
		log.Printf("can't parse csv file: %s", err)
		return
	}

	println(skipped + len(records))
	println(len(records))
	println(skipped)
	spew.Dump(records[38])
}

// ParseCSV parses the CSV content from the io.Reader until EOF, returning
// successfully parsed records and the number of skipped records.
// It validates each record in a streaming manner using validateFunc.
func ParseCSV(r io.Reader, validateFunc func(Record) bool) ([]Record, int, error) {
	dec, err := csvutil.NewDecoder(csvReader(r))
	if err != nil {
		return nil, 0, fmt.Errorf("can't create csv decoder: %w", err)
	}

	var records []Record
	var skipped int
	for {
		var record Record

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
