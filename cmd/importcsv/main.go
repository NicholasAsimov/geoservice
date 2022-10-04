package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/rs/zerolog"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/nicholasasimov/findhotel-assignment/config"
	"github.com/nicholasasimov/findhotel-assignment/csvparse"
	"github.com/nicholasasimov/findhotel-assignment/model"
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

	db, err := gorm.Open(postgres.Open(config.DSN()), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		log.Error().Err(err).Msg("can't connect to database")
		return
	}

	if err = model.MigrateDB(db); err != nil {
		log.Error().Err(err).Msg("can't migrate database")
		return
	}

	file, err := os.Open(config.Config.Importer.Filepath)
	if err != nil {
		log.Error().Err(err).Msg("can't open csv file")
		return
	}
	defer file.Close()

	validate := func(r model.GeoRecord) bool {
		return r.IPAddress.IsValid() && r.City != "" && r.Country != "" && r.CountryCode != ""
	}

	start := time.Now()
	records, skipped, err := csvparse.ParseCSV(file, validate)
	took := time.Since(start)
	if err != nil {
		log.Error().Err(err).Msg("can't parse csv file")
		return
	}

	// note: usually would be implemented in a store package
	result := db.Save(records[0:100])
	if err := result.Error; err != nil {
		log.Error().Err(err).Msg("can't save records in db")
		return
	}

	log.Info().
		Int("records", len(records)+skipped).
		Int("accepted", len(records)).
		Int("skipped", skipped).
		Str("took", took.String()).
		Msg("import finished")
}
