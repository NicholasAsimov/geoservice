package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"net/netip"
	"os"
	"time"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/jszwec/csvutil"
	"github.com/rs/zerolog"

	"github.com/nicholasasimov/geoservice/config"
	"github.com/nicholasasimov/geoservice/model"
)

func main() {
	var records int
	var filename string

	flag.IntVar(&records, "records", 1_000_000, "how many records to generate")
	flag.StringVar(&filename, "file", "fakedata.csv", "file to write csv")
	flag.Parse()

	if err := config.Init(); err != nil {
		fmt.Printf("can't init config: %s\n\n", err)
		config.Usage()
		os.Exit(1)
	}

	log := zerolog.New(os.Stdout).With().Timestamp().Caller().Logger()
	if config.Config.PrettyLog {
		log = log.Output(zerolog.NewConsoleWriter())
	}

	file, err := os.Create(filename)
	if err != nil {
		log.Error().Err(err).Msg("can't create csv file")
		return
	}
	defer file.Close()

	w := csv.NewWriter(file)
	encoder := csvutil.NewEncoder(w)
	defer func() {
		w.Flush()
		if err := w.Error(); err != nil {
			log.Fatal().Err(err).Msg("can't flush file")
		}
	}()

	start := time.Now()
	for i := 0; i < records; i++ {
		record := model.GeoRecord{
			IPAddress:    netip.MustParseAddr(gofakeit.IPv4Address()),
			CountryCode:  gofakeit.CountryAbr(),
			Country:      gofakeit.Country(),
			City:         gofakeit.City(),
			Latitude:     gofakeit.Latitude(),
			Longitude:    gofakeit.Longitude(),
			MysteryValue: gofakeit.Float64Range(0, 1),
		}

		if err := encoder.Encode(record); err != nil {
			log.Error().Err(err).Msg("can't encode csv record")
			return
		}
	}

	log.Info().
		Int("records", records).
		Str("file", filename).
		Str("took", time.Since(start).String()).
		Msg("csv generated")
}
