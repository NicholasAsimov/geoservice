package model

import "net/netip"

type Record struct {
	IPAddress    netip.Addr `csv:"ip_address" gorm:"primaryKey"`
	CountryCode  string     `csv:"country_code"`
	Country      string     `csv:"country"`
	City         string     `csv:"city"`
	Latitude     float64    `csv:"latitude"`
	Longitude    float64    `csv:"longitude"`
	MysteryValue float64    `csv:"mystery_value"`
}

func MigrateDB(db *gorm.DB) error {
	if err := db.AutoMigrate(&model.Record{}); err != nil {
		return err
	}

	return nil
}
