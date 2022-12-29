package store

import (
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"github.com/jackc/pgx/v5"

	"github.com/nicholasasimov/geoservice/resources"
)

func MigrateDB(db *pgx.Conn) error {
	d, err := iofs.New(resources.FS, "sql")
	if err != nil {
		return fmt.Errorf("can't open migrations: %w", err)
	}

	m, err := migrate.NewWithSourceInstance("iofs", d, db.Config().ConnString())
	if err != nil {
		return fmt.Errorf("can't create migrator: %w", err)
	}

	err = m.Up()
	if err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("can't migrate db: %w", err)
	}

	return nil
}
