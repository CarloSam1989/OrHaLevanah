package horizons

import (
	"testing"
	"time"
)

func TestParseEphemerisSelectsFirstConservativeEveningCrescent(t *testing.T) {
	fixture := `header
$$SOE
 2026-Aug-12 16:19:00.000, C, m, 275.1, 1.2, 0.2, 7.1, /T, 172.9,
 2026-Aug-13 16:19:00.000,*,m, 276.7, 7.1, 1.1, 12.1, /T, 167.9,
 2026-Aug-13 16:49:00.000, C, m, 280.2, 3.8, 1.1, 12.2, /T, 167.8,
 2026-Aug-14 16:18, C, m, 285.3, 12.5, 4.2, 23.5, /T, 156.5,
$$EOE`
	loc, err := time.LoadLocation("Asia/Jerusalem")
	if err != nil {
		t.Fatal(err)
	}
	forecast, err := parseEphemeris(fixture, loc)
	if err != nil {
		t.Fatal(err)
	}
	if forecast.PossibleRoshChodesh != "2026-08-13" {
		t.Fatalf("date = %s, want 2026-08-13", forecast.PossibleRoshChodesh)
	}
	if !forecast.EveningCrescent {
		t.Fatal("expected an evening crescent")
	}
}
