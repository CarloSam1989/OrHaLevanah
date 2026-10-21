package sighting

import (
	"strings"
	"testing"
	"time"
)

func TestParseConfirmedIsraelReport(t *testing.T) {
	body := `<h2>Israel-Based Observation</h2><strong>CONFIRMED</strong><p>July 15, 2026</p>`
	report, err := parseConfirmedIsraelReport(body, "https://example.test", time.Date(2026, 7, 15, 22, 0, 0, 0, time.UTC))
	if err != nil {
		t.Fatal(err)
	}
	if report.Date != "2026-07-15" || report.Status != "confirmed" {
		t.Fatalf("unexpected report: %+v", report)
	}
}

func TestParseRejectsUnconfirmedReport(t *testing.T) {
	_, err := parseConfirmedIsraelReport(`<h2>Israel-Based Observation</h2><p>July 15, 2026</p>`, "https://example.test", time.Now())
	if err == nil || !strings.Contains(err.Error(), "no confirmed") {
		t.Fatalf("expected no-confirmation error, got %v", err)
	}
}

func TestParseRejectsGlobalConfirmationAfterUnconfirmedIsraelBlock(t *testing.T) {
	body := `<h2>Israel-Based Observation</h2><strong>NO REPORTS</strong><p>August 13, 2026</p>
		<h2>Global Observation</h2><strong>CONFIRMED</strong><p>August 13, 2026</p>`
	_, err := parseConfirmedIsraelReport(body, "https://example.test", time.Now())
	if err == nil || !strings.Contains(err.Error(), "no confirmed") {
		t.Fatalf("expected global report to be rejected for Israel, got %v", err)
	}
}
