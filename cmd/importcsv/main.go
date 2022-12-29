package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/rs/zerolog"

	"github.com/nicholasasimov/geoservice/config"
	"github.com/nicholasasimov/geoservice/csvparse"
	"github.com/nicholasasimov/geoservice/model"
	"github.com/nicholasasimov/geoservice/store"
)

func main() {
	ctx := context.Background()

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

	db, err := pgx.Connect(ctx, config.DSN())
	if err != nil {
		log.Error().Err(err).Msg("can't connect to database")
		return
	}
	defer db.Close(ctx)

	s := store.New(db)
	if err = store.MigrateDB(db); err != nil {
		log.Error().Err(err).Msg("can't migrate database")
		return
	}

	file, err := os.Open(config.Config.Importer.Filepath)
	if err != nil {
		log.Error().Err(err).Msg("can't open csv file")
		return
	}
	defer file.Close()

	// note: since validation (business) logic is separated from parsing it can
	// be, for example, be configurable by external configuration.
	validate := func(r model.GeoRecord) bool {
		return r.IPAddress.IsValid() && r.City != "" && r.Country != "" && r.CountryCode != ""
	}

	log.Info().Str("file", file.Name()).Msg("parsing file")
	parseStart := time.Now()

	records, skipped, err := csvparse.ParseCSV(file, validate)
	if err != nil {
		log.Error().Err(err).Msg("can't parse csv file")
		return
	}

	parseTook := time.Since(parseStart)

	log.Info().Int("records", len(records)).Msg("persisting to db")
	dbStart := time.Now()

	err = s.UpsertRecords(ctx, records)
	if err != nil {
		log.Error().Err(err).Msg("can't save records in db")
		return
	}

	dbTook := time.Since(dbStart)

	log.Info().
		Int("records", len(records)+skipped).
		Int("accepted", len(records)).
		Int("skipped", skipped).
		Str("parse_took", parseTook.String()).
		Str("db_took", dbTook.String()).
		Msg("import finished")
}
