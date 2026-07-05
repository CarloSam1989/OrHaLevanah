package biblical

import (
	"testing"
	"time"
)

func TestEstimateBiblicalYearStartStaysNearSpringEquinox(t *testing.T) {
	loc, err := time.LoadLocation("Asia/Jerusalem")
	if err != nil {
		t.Fatal(err)
	}

	now := time.Date(2026, time.July, 5, 12, 0, 0, 0, loc)
	start := EstimateBiblicalYearStart(now)
	t.Logf("inicio estimado de Aviv: %s", start.Format("2006-01-02"))
	earliest := time.Date(2026, time.March, 20, 0, 0, 0, 0, loc)
	latest := time.Date(2026, time.May, 5, 0, 0, 0, 0, loc)

	if start.Before(earliest) || start.After(latest) {
		t.Fatalf("el inicio de Aviv quedó fuera de la ventana anual: %s", start.Format("2006-01-02"))
	}
}
