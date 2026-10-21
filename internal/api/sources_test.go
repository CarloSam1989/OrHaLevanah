package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"or-halevanah/internal/biblical"
	"or-halevanah/internal/models"
)

func sourceFixtures(t *testing.T, reports, forecast string) {
	t.Helper()
	dir := t.TempDir()
	for env, entry := range map[string]struct{ name, content string }{
		"MOON_SIGHTING_CACHE_PATH": {"reports.json", reports},
		"JPL_HORIZONS_CACHE_PATH":  {"forecast.json", forecast},
	} {
		path := filepath.Join(dir, entry.name)
		if err := os.WriteFile(path, []byte(entry.content), 0600); err != nil {
			t.Fatal(err)
		}
		t.Setenv(env, path)
	}
}

func sourcesResult(t *testing.T, date string) struct {
	Date    string        `json:"date"`
	IsHead  bool          `json:"is_month_start"`
	Sources []monthSource `json:"sources"`
} {
	t.Helper()
	w := httptest.NewRecorder()
	MonthSourcesHandler(w, httptest.NewRequest("GET", "/?date="+date, nil))
	if w.Code != 200 {
		t.Fatalf("status %d: %s", w.Code, w.Body.String())
	}
	var result struct {
		Data struct {
			Date    string        `json:"date"`
			IsHead  bool          `json:"is_month_start"`
			Sources []monthSource `json:"sources"`
		} `json:"data"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &result); err != nil {
		t.Fatal(err)
	}
	return result.Data
}

func TestSourcesPeriodAndIsraelScope(t *testing.T) {
	sourceFixtures(t, `{"reports":[
	{"date":"2026-08-14","country":"Israel","status":"confirmed","source_name":"Eligible","source_url":"https://example.test/report"},
	{"date":"2026-08-14","country":"Israel","status":"confirmed","source_name":"Eligible","source_url":"https://example.test/report"},
	{"date":"2026-08-14","country":"Colombia","status":"confirmed","source_name":"Other country"},
	{"date":"2026-07-15","country":"Israel","status":"confirmed","source_name":"Other period"},
	{"date":"2026-08-14","country":"Israel","status":"pending","source_name":"Pending"}
	]}`, `{"forecast":{}}`)
	got := sourcesResult(t, "2026-08-14&country=Colombia&is_month_start=false&user_id=another")
	if !got.IsHead || len(got.Sources) != 1 || got.Sources[0].Name != "Eligible" {
		t.Fatalf("wrong scope: %+v", got)
	}
	ordinary := sourcesResult(t, "2026-08-15&is_month_start=true")
	if ordinary.IsHead || len(ordinary.Sources) != 0 {
		t.Fatalf("ordinary date: %+v", ordinary)
	}
}

func TestEstimatedHeadWithoutSources(t *testing.T) {
	sourceFixtures(t, `{"reports":[]}`, `{"forecast":{}}`)
	loc, _ := time.LoadLocation("Asia/Jerusalem")
	start := biblical.FindEstimatedMonthStart(time.Date(2026, 9, 28, 23, 0, 0, 0, loc))
	got := sourcesResult(t, start.Format("2006-01-02"))
	if !got.IsHead || got.Sources == nil || len(got.Sources) != 0 {
		t.Fatalf("expected empty head: %+v", got)
	}
}

func TestForecastBelongsToPreviousMonth(t *testing.T) {
	for _, anchor := range []string{"2026-08-14", "2025-08-14"} {
		t.Run(anchor, func(t *testing.T) {
			sourceFixtures(t, `{"reports":[{"date":"2026-08-14","country":"Israel","status":"confirmed","source_name":"Report"}]}`,
				`{"forecast":{"possible_rosh_chodesh":"2026-09-12","based_on_month_start":"`+anchor+`","source":"NASA","source_url":"https://example.test/forecast"}}`)
			got := sourcesResult(t, "2026-09-12")
			if anchor == "2026-08-14" {
				if !got.IsHead || len(got.Sources) != 1 || got.Sources[0].Status != "estimated" {
					t.Fatalf("forecast: %+v", got)
				}
			} else if len(got.Sources) != 0 {
				t.Fatalf("wrong period forecast: %+v", got)
			}
		})
	}
}

func TestSourceValidationAndErrors(t *testing.T) {
	sourceFixtures(t, `{"reports":[]}`, `{"forecast":{}}`)
	for _, date := range []string{"", "2026-02-30", "1899-01-01", "2201-01-01", "2026-9-1"} {
		w := httptest.NewRecorder()
		MonthSourcesHandler(w, httptest.NewRequest("GET", "/?date="+date, nil))
		if w.Code != 400 {
			t.Errorf("%s returned %d", date, w.Code)
		}
	}
	w := httptest.NewRecorder()
	MonthSourcesHandler(w, httptest.NewRequest("POST", "/?date=2026-08-14", nil))
	if w.Code != http.StatusMethodNotAllowed {
		t.Fatal(w.Code)
	}
	if err := os.WriteFile(os.Getenv("MOON_SIGHTING_CACHE_PATH"), []byte("invalid"), 0600); err != nil {
		t.Fatal(err)
	}
	w = httptest.NewRecorder()
	MonthSourcesHandler(w, httptest.NewRequest("GET", "/?date=2026-08-14", nil))
	if w.Code != 500 {
		t.Fatalf("corrupt storage must not appear empty: %d", w.Code)
	}
	for _, raw := range []string{"javascript:alert(1)", "//evil.test", "https://user:secret@example.test"} {
		if sourceLink(raw) != "" {
			t.Errorf("unsafe link: %s", raw)
		}
	}
}

func TestTodayUsesJerusalemServerDate(t *testing.T) {
	sourceFixtures(t, `{"reports":[]}`, `{"forecast":{}}`)
	loc, _ := time.LoadLocation("Asia/Jerusalem")
	before := time.Now().In(loc).Format("2006-01-02")
	w := httptest.NewRecorder()
	BiblicalJerusalemTodayHandler(w, httptest.NewRequest("GET", "/?date=1900-01-01&timezone=America/Bogota", nil))
	var result struct {
		Data models.BiblicalTodayResponse `json:"data"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &result); err != nil {
		t.Fatal(err)
	}
	after := time.Now().In(loc).Format("2006-01-02")
	if result.Data.CivilDate != before && result.Data.CivilDate != after {
		t.Fatalf("wrong server date: %+v", result.Data)
	}
	if result.Data.Location.Timezone != "Asia/Jerusalem" {
		t.Fatal(result.Data.Location)
	}
}

func TestExistingAvivRegistration(t *testing.T) {
	dir := t.TempDir()
	reportsPath := filepath.Join(dir, "aviv_reports.json")
	statusPath := filepath.Join(dir, "aviv_status.json")
	t.Setenv("AVIV_REPORTS_PATH", reportsPath)
	t.Setenv("AVIV_STATUS_PATH", statusPath)
	for path, value := range map[string]string{reportsPath: "[]", statusPath: "{}"} {
		if err := os.WriteFile(path, []byte(value), 0600); err != nil {
			t.Fatal(err)
		}
	}
	for _, body := range []string{
		`{"user_name":"Test","confirmed":false,"notes":"Pending"}`,
		`{"user_name":"Test","confirmed":true,"evidence_url":"https://example.test/aviv"}`,
	} {
		w := httptest.NewRecorder()
		AvivReportHandler(w, httptest.NewRequest("POST", "/api/aviv/report", strings.NewReader(body)))
		if w.Code != 200 {
			t.Fatalf("registration: %d %s", w.Code, w.Body.String())
		}
	}
	w := httptest.NewRecorder()
	AvivReportHandler(w, httptest.NewRequest("POST", "/api/aviv/report", strings.NewReader(`{"confirmed":true}`)))
	if w.Code != 400 {
		t.Fatalf("missing name: %d", w.Code)
	}
}

func TestSourcesRecognizeBoundaryAlreadyExposedByMonth(t *testing.T) {
	sourceFixtures(t, `{"reports":[{"date":"2026-08-14","country":"Israel","status":"confirmed","source_name":"Report"}]}`, `{"forecast":{}}`)
	loc, _ := time.LoadLocation("Asia/Jerusalem")
	now := time.Date(2026, 9, 28, 12, 0, 0, 0, loc)
	start := biblical.FindEstimatedMonthStart(now)
	if !currentMonthBoundary(start, now) {
		t.Fatal("current calculated month must remain consultable")
	}
	next := biblical.FindNextEstimatedMonthStart(start)
	if !currentMonthBoundary(next, now) {
		t.Fatal("next calculated boundary must remain consultable")
	}
	if currentMonthBoundary(start.AddDate(0, 0, 5), now) {
		t.Fatal("ordinary date must not become a boundary")
	}
}
