package solar

import (
	"testing"
	"time"
)

func TestApproxSunsetJerusalemAppliesDaylightSavingTime(t *testing.T) {
	loc, err := time.LoadLocation("Asia/Jerusalem")
	if err != nil {
		t.Fatal(err)
	}
	summer := ApproxSunsetJerusalem(time.Date(2026, time.July, 15, 12, 0, 0, 0, loc))
	if summer.Hour() != 19 {
		t.Fatalf("summer sunset hour = %d, want 19 with IDT", summer.Hour())
	}
	winter := ApproxSunsetJerusalem(time.Date(2026, time.January, 15, 12, 0, 0, 0, loc))
	if winter.Hour() != 16 {
		t.Fatalf("winter sunset hour = %d, want 16 without DST", winter.Hour())
	}
}
