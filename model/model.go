package model

import (
	"net/netip"
)

// note: depending on the complexity of the data types the DB model would be
// separate from CSV record model. In this case reusing it makes things simpler.
type GeoRecord struct {
	IPAddress    netip.Addr `csv:"ip_address" json:"ip_address" db:"ip_address"`
	CountryCode  string     `csv:"country_code" json:"country_code" db:"country_code"`
	Country      string     `csv:"country" json:"country" db:"country"`
	City         string     `csv:"city" json:"city" db:"city"`
	Latitude     float64    `csv:"latitude" json:"latitude" db:"latitude"`
	Longitude    float64    `csv:"longitude" json:"longitude" db:"longitude"`
	MysteryValue float64    `csv:"mystery_value" json:"mystery_value" db:"mystery_value"`
}
