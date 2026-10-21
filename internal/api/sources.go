package api

import (
	"encoding/json"
	"net/http"
	"net/url"
	"os"
	"sort"
	"strings"
	"time"

	"or-halevanah/internal/biblical"
	"or-halevanah/internal/horizons"
	"or-halevanah/internal/location"
	"or-halevanah/internal/models"
	"or-halevanah/internal/sighting"
)

// Public calendar evidence only: never includes community members' Aviv reports.
type monthSource struct {
	Name   string `json:"name"`
	URL    string `json:"url,omitempty"`
	Status string `json:"status"`
	Method string `json:"method,omitempty"`
	Title string `json:"title,omitempty"`
	Date string `json:"date,omitempty"`
}

func sourceLink(raw string) string {
	u, err := url.Parse(raw)
	if err != nil || u.Host == "" || (u.Scheme != "https" && u.Scheme != "http") || u.User != nil {
		return ""
	}
	return u.String()
}

func MonthSourcesHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-Control", "no-store")
	if r.Method != http.MethodGet {
		w.Header().Set("Allow", "GET")
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}
	loc, err := time.LoadLocation(location.GetJerusalem().Timezone)
	if err != nil {
		http.Error(w, "Zona horaria no disponible", http.StatusInternalServerError)
		return
	}
	date, err := time.ParseInLocation("2006-01-02", r.URL.Query().Get("date"), loc)
	if err != nil || date.Year() < 1900 || date.Year() > 2200 {
		http.Error(w, "Fecha inválida: use YYYY-MM-DD entre 1900 y 2200", http.StatusBadRequest)
		return
	}
	reports, err := sighting.LoadReports()
	if err != nil && !os.IsNotExist(err) {
		http.Error(w, "No se pudieron consultar las fuentes", http.StatusInternalServerError)
		return
	}
	key := date.Format("2006-01-02")
	sources := make([]monthSource, 0)
	seen := make(map[monthSource]bool)
	isHead := false
	for _, report := range reports {
		if report.Date != key || report.Status != "confirmed" || !strings.EqualFold(report.Country, "Israel") {
			continue
		}
		isHead = true
		source := monthSource{Name: report.SourceName, URL: sourceLink(report.SourceURL), Status: report.Status, Method: report.ConfirmationMethod, Date: report.Date}
		if report.ConfirmationMethod == "day_30_declaration" { source.Status = "declared" }
		if strings.TrimSpace(source.Name) == "" && source.URL == "" {
			continue
		}
		if !seen[source] {
			sources = append(sources, source)
			seen[source] = true
		}
	}
	// A calendar cell describes the biblical day beginning at its sunset.
	if !isHead {
		start := biblical.FindEstimatedMonthStart(time.Date(date.Year(), date.Month(), date.Day(), 23, 0, 0, 0, loc))
		isHead = start.Format("2006-01-02") == key
		if !isHead {
			isHead = currentMonthBoundary(date, time.Now().In(loc))
		}
	}
	forecast, forecastErr := horizons.LoadForecast()
	if forecastErr != nil && !os.IsNotExist(forecastErr) {
		http.Error(w, "No se pudieron consultar las fuentes", http.StatusInternalServerError)
		return
	}
	if forecastErr == nil && forecast.PossibleRoshChodesh == key {
		previous := biblical.FindEstimatedMonthStart(date) // before sunset
		if forecast.BasedOnMonthStart == previous.Format("2006-01-02") {
			isHead = true
			if forecast.Source != "" || sourceLink(forecast.SourceURL) != "" {
				sources = append(sources, monthSource{Name: forecast.Source, URL: sourceLink(forecast.SourceURL), Status: "estimated", Date: forecast.PossibleRoshChodesh})
			}
		}
	}
	evidence, searchState := sighting.EvidenceFor(date)
	if isHead {
		for _, item := range evidence {
			status := item.Status
			if item.Date != key && (status == "confirmed" || status == "declared") { status = "related" }
			sources = append(sources, monthSource{Name: "Devorah's Date Tree", URL: sourceLink(item.URL), Status: status, Method: item.Method, Title: item.Title, Date: item.Date})
		}
		// Also retain nearby cached evidence, but never label it a confirmation
		// of the selected day. A different observation date stays explicit.
		for _, report := range reports {
			if report.Date == key || !strings.EqualFold(report.Country, "Israel") { continue }
			if report.Date < date.AddDate(0,0,-7).Format("2006-01-02") || report.Date > date.AddDate(0,0,7).Format("2006-01-02") { continue }
			sources = append(sources, monthSource{Name: report.SourceName, URL: sourceLink(report.SourceURL), Status: "related", Method: report.ConfirmationMethod, Date: report.Date})
		}
		if forecastErr == nil && forecast.PossibleRoshChodesh != key && forecast.Source != "" {
			if predicted, err := time.ParseInLocation("2006-01-02", forecast.PossibleRoshChodesh, loc); err == nil && predicted.After(date.AddDate(0,0,-3)) && predicted.Before(date.AddDate(0,0,3)) {
				previous := biblical.FindEstimatedMonthStart(predicted)
				if forecast.BasedOnMonthStart == previous.Format("2006-01-02") {
					sources = append(sources, monthSource{Name: forecast.Source, URL: sourceLink(forecast.SourceURL), Status: "estimated", Date: forecast.PossibleRoshChodesh})
				}
			}
		}
	}
	sources = uniqueSources(sources)
	json.NewEncoder(w).Encode(models.ApiResponse{Success: true, Data: map[string]interface{}{
		"date": key, "is_month_start": isHead, "sources": sources, "search_state": searchState,
	}})
}

func uniqueSources(sources []monthSource) []monthSource {
	rank := map[string]int{"confirmed":0, "declared":1, "related":2, "estimated":3, "pending":4}
	sort.SliceStable(sources,func(i,j int) bool { return rank[sources[i].Status] < rank[sources[j].Status] })
	result := make([]monthSource,0,len(sources))
	positions := map[string]int{}
	for _, source := range sources {
		if source.Name == "" && source.URL == "" { continue }
		key := source.URL
		if key == "" { key = source.Name+"|"+source.Date }
		if i, found := positions[key]; found {
			if result[i].Title == "" { result[i].Title = source.Title }
			continue
		}
		positions[key] = len(result)
		result = append(result,source)
	}
	return result
}

// /month may expose a calculated boundary after its preceding confirmation
// ages out. Validate that same boundary without changing the calendar algorithm.
func currentMonthBoundary(date, now time.Time) bool {
	start := biblical.FindEstimatedMonthStart(now)
	if date.Format("2006-01-02") == start.Format("2006-01-02") {
		return true
	}
	return date.Format("2006-01-02") == biblical.FindNextEstimatedMonthStart(start).Format("2006-01-02")
}
