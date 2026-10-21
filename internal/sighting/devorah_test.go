package sighting

import (
	"testing"
	"time"
)

func TestDevorahReportLinksRejectsUpcomingAnnouncements(t *testing.T) {
	body := `<a href="https://mailchi.mp/devorahsdatetree/upcoming-moon-observation-7th-biblical-month">Upcoming</a>
		<a href="https://mailchi.mp/devorahsdatetree/7th-biblical-month-crescent-moon-report?e=abc">Report</a>`
	links := devorahReportLinks(body)
	if len(links) != 1 || links[0] != "https://mailchi.mp/devorahsdatetree/7th-biblical-month-crescent-moon-report" {
		t.Fatalf("unexpected report links: %#v", links)
	}
}

func TestParseDevorahIsraelEyewitnessReport(t *testing.T) {
	body := `<h1>The 6th Biblical Month has begun!</h1><p>The crescent moon was sighted on 14 August 2026 from Haifa.</p><p>Observers looked for the moon from Israel.</p>`
	report, err := parseDevorahMonthReport(body, "https://mailchi.mp/devorahsdatetree/report", time.Now())
	if err != nil {
		t.Fatal(err)
	}
	if report.Date != "2026-08-14" || report.BiblicalMonth != 6 || report.ConfirmationMethod != "israel_eyewitness" {
		t.Fatalf("unexpected report: %+v", report)
	}
}

func TestParseDevorahRejectsNegativeIsraelReport(t *testing.T) {
	body := `<h1>6th Biblical Month Report</h1><p>We aren't aware of the crescent moon being spotted from Israel this evening, 13 August 2026.</p>`
	if _, err := parseDevorahMonthReport(body, "https://example.test", time.Now()); err == nil {
		t.Fatal("expected negative report to be rejected")
	}
}

func TestParseDevorahDayThirtyDeclaration(t *testing.T) {
	body := `<h1>3rd Biblical Month Report</h1><p>As this was the end of the 29th day, the Biblical Month will begin on 18 May 2026 at sunset by default.</p>`
	report, err := parseDevorahMonthReport(body, "https://mailchi.mp/devorahsdatetree/report", time.Now())
	if err != nil {
		t.Fatal(err)
	}
	if report.Date != "2026-05-18" || report.ConfirmationMethod != "day_30_declaration" {
		t.Fatalf("unexpected report: %+v", report)
	}
}
