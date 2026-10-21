package biblical

import (
	"path/filepath"
	"testing"
	"time"
)

func TestConfirmedIsraelSightingOverridesEstimate(t *testing.T) {
	cache, err := filepath.Abs(filepath.Join("..", "..", "data", "moon_sightings.json"))
	if err != nil {
		t.Fatal(err)
	}
	t.Setenv("MOON_SIGHTING_CACHE_PATH", cache)

	loc, err := time.LoadLocation("Asia/Jerusalem")
	if err != nil {
		t.Fatal(err)
	}
	now := time.Date(2026, time.July, 18, 12, 0, 0, 0, loc)
	start := FindEstimatedMonthStart(now)
	if got, want := start.Format("2006-01-02"), "2026-07-15"; got != want {
		t.Fatalf("month start = %s, want confirmed sighting %s", got, want)
	}
	if got, want := CalculateBiblicalDay(start, now, false), 3; got != want {
		t.Fatalf("biblical day = %d, want %d", got, want)
	}
}

func TestAugustIsraelSightingStartsMonthOnAugust14(t *testing.T) {
	cache, err := filepath.Abs(filepath.Join("..", "..", "data", "moon_sightings.json"))
	if err != nil {
		t.Fatal(err)
	}
	t.Setenv("MOON_SIGHTING_CACHE_PATH", cache)

	loc, err := time.LoadLocation("Asia/Jerusalem")
	if err != nil {
		t.Fatal(err)
	}
	now := time.Date(2026, time.August, 15, 12, 0, 0, 0, loc)
	start := FindEstimatedMonthStart(now)
	if got, want := start.Format("2006-01-02"), "2026-08-14"; got != want {
		t.Fatalf("month start = %s, want confirmed Israel sighting %s", got, want)
	}
}
