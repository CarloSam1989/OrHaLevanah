package biblical

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestConfirmedAvivMakesJulyTheFifthMonth(t *testing.T) {
	root, err := filepath.Abs(filepath.Join("..", ".."))
	if err != nil {
		t.Fatal(err)
	}
	t.Setenv("AVIV_STATUS_PATH", filepath.Join(root, "data", "aviv_status.json"))
	t.Setenv("MOON_SIGHTING_CACHE_PATH", filepath.Join(root, "data", "moon_sightings.json"))
	t.Setenv("JPL_HORIZONS_CACHE_PATH", filepath.Join(root, "data", "horizons_jerusalem.json"))
	if _, err := os.Stat(filepath.Join(root, "data", "aviv_status.json")); err != nil {
		t.Fatal(err)
	}
	loc, err := time.LoadLocation("Asia/Jerusalem")
	if err != nil {
		t.Fatal(err)
	}

	now := time.Date(2026, time.July, 18, 12, 0, 0, 0, loc)
	if got, want := EstimateBiblicalYearStart(now).Format("2006-01-02"), "2026-03-20"; got != want {
		t.Fatalf("year start = %s, want %s", got, want)
	}
	if got, want := GetCurrentBiblicalMonthAt(now), 5; got != want {
		t.Fatalf("biblical month = %d, want %d", got, want)
	}
}

func TestMonthChangesAtJerusalemSunsetNotMidnight(t *testing.T) {
	root, err := filepath.Abs(filepath.Join("..", ".."))
	if err != nil {
		t.Fatal(err)
	}
	t.Setenv("AVIV_STATUS_PATH", filepath.Join(root, "data", "aviv_status.json"))
	t.Setenv("MOON_SIGHTING_CACHE_PATH", filepath.Join(root, "data", "moon_sightings.json"))
	t.Setenv("JPL_HORIZONS_CACHE_PATH", filepath.Join(root, "data", "horizons_jerusalem.json"))
	loc, err := time.LoadLocation("Asia/Jerusalem")
	if err != nil {
		t.Fatal(err)
	}
	beforeSunset := time.Date(2026, time.July, 15, 12, 0, 0, 0, loc)
	afterSunset := time.Date(2026, time.July, 15, 22, 0, 0, 0, loc)
	if got := GetCurrentBiblicalMonthAt(beforeSunset); got != 4 {
		t.Fatalf("month before sunset = %d, want 4", got)
	}
	if got := GetCurrentBiblicalMonthAt(afterSunset); got != 5 {
		t.Fatalf("month after sunset = %d, want 5", got)
	}
}
