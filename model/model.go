package model

import (
	"net/netip"

	"gorm.io/gorm"
)

// note: depending on the complexity of the data types the DB model would be
// separate from CSV record model. In this case reusing it makes things simpler.
type GeoRecord struct {
	IPAddress    netip.Addr `csv:"ip_address" gorm:"primaryKey;serializer:json"`
	CountryCode  string     `csv:"country_code"`
	Country      string     `csv:"country"`
	City         string     `csv:"city"`
	Latitude     float64    `csv:"latitude"`
	Longitude    float64    `csv:"longitude"`
	MysteryValue float64    `csv:"mystery_value"`
}

func MigrateDB(db *gorm.DB) error {
	if err := db.AutoMigrate(&GeoRecord{}); err != nil {
		return err
	}

	return nil
}
