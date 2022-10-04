package model

import (
	"net/netip"
)

// note: depending on the complexity of the data types the DB model would be
// separate from CSV record model. In this case reusing it makes things simpler.
type GeoRecord struct {
	IPAddress    netip.Addr `csv:"ip_address"`
	CountryCode  string     `csv:"country_code"`
	Country      string     `csv:"country"`
	City         string     `csv:"city"`
	Latitude     float64    `csv:"latitude"`
	Longitude    float64    `csv:"longitude"`
	MysteryValue float64    `csv:"mystery_value"`
}
