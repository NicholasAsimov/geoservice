package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"runtime/pprof"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/rs/zerolog"

	"github.com/nicholasasimov/findhotel-assignment/config"
	"github.com/nicholasasimov/findhotel-assignment/csvparse"
	"github.com/nicholasasimov/findhotel-assignment/model"
	"github.com/nicholasasimov/findhotel-assignment/store"
)

func main() {
	flag.Usage = config.Usage
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

	db, err := pgx.Connect(context.Background(), config.DSN())
	if err != nil {
		log.Error().Err(err).Msg("can't connect to database")
		return
	}
	defer db.Close(context.Background())

	// if err = model.MigrateDB(db); err != nil {
	// 	log.Error().Err(err).Msg("can't migrate database")
	// 	return
	// }

	file, err := os.Open(config.Config.Importer.Filepath)
	if err != nil {
		log.Error().Err(err).Msg("can't open csv file")
		return
	}
	defer file.Close()

	validate := func(r model.GeoRecord) bool {
		return r.IPAddress.IsValid() && r.City != "" && r.Country != "" && r.CountryCode != ""
	}

	parseStart := time.Now()
	records, skipped, err := csvparse.ParseCSV(file, validate)
	parseTook := time.Since(parseStart)
	if err != nil {
		log.Error().Err(err).Msg("can't parse csv file")
		return
	}

	//  TODO move to store
	dbStart := time.Now()
	copyCount, err := db.CopyFrom(
		context.Background(),
		pgx.Identifier{"geo_records"},
		[]string{"ip_address", "country_code", "country", "city", "latitude", "longitude", "mystery_value"},
		store.CopyFromRecords(records),
	)
	dbTook := time.Since(dbStart)
	if err != nil {
		log.Error().Err(err).Msg("can't save records in db")
		return
	}

	log.Info().
		Int("records", len(records)+skipped).
		Int("accepted", len(records)).
		Int("skipped", skipped).
		Int64("copy_count", copyCount).
		Str("parse_took", parseTook.String()).
		Str("db_took", dbTook.String()).
		Msg("import finished")

	f, _ := os.Create("heap_profile")
	pprof.WriteHeapProfile(f)
	f.Close()
}
