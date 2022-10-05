package store

import (
	"testing"
)

func TestBuildSetSQL(t *testing.T) {
	in := []string{"a", "bb", "ccc", "dddd"}
	want := `a=EXCLUDED.a, bb=EXCLUDED.bb, ccc=EXCLUDED.ccc, dddd=EXCLUDED.dddd`
	got := buildSetSQL(in)

	if want != got {
		t.Errorf("\n expected \t%q\n got \t\t%q", want, got)
	}
}
