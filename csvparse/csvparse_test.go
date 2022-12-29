package csvparse_test

import (
	"net/netip"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/nicholasasimov/geoservice/csvparse"
	"github.com/nicholasasimov/geoservice/model"
)

func TestParseCSV(t *testing.T) {
	tests := map[string]struct {
		in       string
		want     []model.GeoRecord
		validate func(model.GeoRecord) bool
		skipped  int
	}{
		"duplicate ips first record is used": {
			in: `ip_address,country_code,country,city,latitude,longitude,mystery_value
			8.8.8.8,SI,Nepal,DuBuquemouth,10,20,30
			8.8.8.8,CZ,Nicaragua,New Neva,40,50,60`,
			want: []model.GeoRecord{
				{
					IPAddress:    netip.MustParseAddr("8.8.8.8"),
					CountryCode:  "SI",
					Country:      "Nepal",
					City:         "DuBuquemouth",
					Latitude:     10,
					Longitude:    20,
					MysteryValue: 30,
				},
			},
			validate: func(r model.GeoRecord) bool { return true },
			skipped:  1,
		},
		"full test with validation": {
			in: `ip_address,country_code,country,city,latitude,longitude,mystery_value
			200.106.141.15,SI,Nepal,DuBuquemouth,-84.87503094689836,7.206435933364332,782301134
			160.103.7.140,CZ,Nicaragua,New Neva,-68.31023296602508,-37.62435199624531,7301823115
			70.95.73.73,TL,Saudi Arabia,Gradymouth,-49.16675918861615,-86.05920084416894,2559997162
			,PY,Falkland Islands (Malvinas),,75.41685191518815,-144.6943217219469,0`,
			want: []model.GeoRecord{
				{
					IPAddress:    netip.MustParseAddr("200.106.141.15"),
					CountryCode:  "SI",
					Country:      "Nepal",
					City:         "DuBuquemouth",
					Latitude:     -84.87503094689836,
					Longitude:    7.206435933364332,
					MysteryValue: 782301134,
				},
				{
					IPAddress:    netip.MustParseAddr("160.103.7.140"),
					CountryCode:  "CZ",
					Country:      "Nicaragua",
					City:         "New Neva",
					Latitude:     -68.31023296602508,
					Longitude:    -37.62435199624531,
					MysteryValue: 7301823115,
				},
				{
					IPAddress:    netip.MustParseAddr("70.95.73.73"),
					CountryCode:  "TL",
					Country:      "Saudi Arabia",
					City:         "Gradymouth",
					Latitude:     -49.16675918861615,
					Longitude:    -86.05920084416894,
					MysteryValue: 2559997162,
				},
			},
			validate: func(r model.GeoRecord) bool { return r.IPAddress.IsValid() },
			skipped:  1,
		},
	}

	opts := []cmp.Option{
		cmp.Comparer(func(a, b netip.Addr) bool { return a.String() == b.String() }),
	}

	for name, tt := range tests {
		tt := tt

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got, gotSkipped, err := csvparse.ParseCSV(strings.NewReader(tt.in), tt.validate)
			if err != nil {
				t.Errorf("unexpected err: %s", err)
			}

			if diff := cmp.Diff(tt.want, got, opts...); diff != "" {
				t.Errorf("mismatch (-want +got):\n%s", diff)
			}

			if tt.skipped != gotSkipped {
				t.Errorf("expected %d skipped, got %d", tt.skipped, gotSkipped)
			}
		})
	}
}
