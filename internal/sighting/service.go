package sighting

import (
	"encoding/json"
	"errors"
	"fmt"
	"html"
	"io"
	"net/http"
	"os"
	"regexp"
	"sort"
	"strings"
	"sync"
	"time"
)

const (
	defaultSourceURL = "https://renewedmoon.com/"
	defaultCachePath = "data/moon_sightings.json"
)

type Report struct {
	Date       string    `json:"date"`
	Country    string    `json:"country"`
	Location   string    `json:"location,omitempty"`
	Status     string    `json:"status"`
	SourceName string    `json:"source_name"`
	SourceURL  string    `json:"source_url"`
	CheckedAt  time.Time `json:"checked_at"`
}

type cacheFile struct {
	Reports []Report `json:"reports"`
}

var (
	mu          sync.RWMutex
	lastAttempt time.Time
	refreshTTL  = 15 * time.Minute
	datePattern = regexp.MustCompile(`(?is)Israel-Based Observation.*?CONFIRMED.*?(January|February|March|April|May|June|July|August|September|October|November|December)\s+(\d{1,2}),\s+(\d{4})`)
	tagPattern  = regexp.MustCompile(`<[^>]+>`)
)

func sourceURL() string {
	if value := strings.TrimSpace(os.Getenv("MOON_SIGHTING_SOURCE_URL")); value != "" {
		return value
	}
	return defaultSourceURL
}

func cachePath() string {
	if value := strings.TrimSpace(os.Getenv("MOON_SIGHTING_CACHE_PATH")); value != "" {
		return value
	}
	return defaultCachePath
}

func LoadReports() ([]Report, error) {
	mu.RLock()
	defer mu.RUnlock()
	return loadReportsUnlocked()
}

func loadReportsUnlocked() ([]Report, error) {
	content, err := os.ReadFile(cachePath())
	if err != nil {
		return nil, err
	}

	var cache cacheFile
	if err := json.Unmarshal(content, &cache); err != nil {
		return nil, err
	}
	return cache.Reports, nil
}

func RefreshIfStale() error {
	mu.Lock()
	if time.Since(lastAttempt) < refreshTTL {
		mu.Unlock()
		return nil
	}
	lastAttempt = time.Now()
	mu.Unlock()

	client := &http.Client{Timeout: 12 * time.Second}
	req, err := http.NewRequest(http.MethodGet, sourceURL(), nil)
	if err != nil {
		return err
	}
	req.Header.Set("User-Agent", "Or-HaLevanah/1.0 moon-sighting-checker")

	response, err := client.Do(req)
	if err != nil {
		return err
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("moon sighting source returned HTTP %d", response.StatusCode)
	}

	body, err := io.ReadAll(io.LimitReader(response.Body, 2<<20))
	if err != nil {
		return err
	}
	report, err := parseConfirmedIsraelReport(string(body), sourceURL(), time.Now().UTC())
	if err != nil {
		return err
	}
	return saveReport(report)
}

// StartAutoRefresh keeps checking the external source while the server runs.
// The cached confirmation remains authoritative if a refresh fails.
func StartAutoRefresh(stop <-chan struct{}) {
	go func() {
		_ = RefreshIfStale()
		ticker := time.NewTicker(refreshTTL)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				_ = RefreshIfStale()
			case <-stop:
				return
			}
		}
	}()
}

func parseConfirmedIsraelReport(body, url string, checkedAt time.Time) (Report, error) {
	plain := html.UnescapeString(tagPattern.ReplaceAllString(body, " "))
	plain = strings.Join(strings.Fields(plain), " ")
	match := datePattern.FindStringSubmatch(plain)
	if len(match) != 4 {
		return Report{}, errors.New("no confirmed Israel sighting found in source")
	}

	date, err := time.Parse("January 2, 2006", match[1]+" "+match[2]+", "+match[3])
	if err != nil {
		return Report{}, err
	}
	return Report{
		Date:       date.Format("2006-01-02"),
		Country:    "Israel",
		Status:     "confirmed",
		SourceName: "Renewed Moon",
		SourceURL:  url,
		CheckedAt:  checkedAt,
	}, nil
}

func saveReport(report Report) error {
	mu.Lock()
	defer mu.Unlock()

	reports, err := loadReportsUnlocked()
	if err != nil && !os.IsNotExist(err) {
		return err
	}
	found := false
	for i := range reports {
		if reports[i].Date == report.Date && reports[i].Country == report.Country {
			if report.Location == "" {
				report.Location = reports[i].Location
			}
			reports[i] = report
			found = true
			break
		}
	}
	if !found {
		reports = append(reports, report)
	}
	sort.Slice(reports, func(i, j int) bool { return reports[i].Date < reports[j].Date })

	content, err := json.MarshalIndent(cacheFile{Reports: reports}, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(cachePath(), append(content, '\n'), 0644)
}

func LatestConfirmedOnOrBefore(day time.Time, maxAgeDays int) (time.Time, Report, bool) {
	reports, err := LoadReports()
	if err != nil {
		return time.Time{}, Report{}, false
	}
	target := time.Date(day.Year(), day.Month(), day.Day(), 0, 0, 0, 0, day.Location())
	var best time.Time
	var selected Report
	for _, report := range reports {
		if report.Status != "confirmed" || !strings.EqualFold(report.Country, "Israel") {
			continue
		}
		parsed, err := time.ParseInLocation("2006-01-02", report.Date, day.Location())
		if err != nil || parsed.After(target) || target.Sub(parsed) > time.Duration(maxAgeDays)*24*time.Hour {
			continue
		}
		if best.IsZero() || parsed.After(best) {
			best, selected = parsed, report
		}
	}
	return best, selected, !best.IsZero()
}

func FirstConfirmedBetween(start, end time.Time) (time.Time, Report, bool) {
	reports, err := LoadReports()
	if err != nil {
		return time.Time{}, Report{}, false
	}
	start = time.Date(start.Year(), start.Month(), start.Day(), 0, 0, 0, 0, start.Location())
	end = time.Date(end.Year(), end.Month(), end.Day(), 0, 0, 0, 0, end.Location())
	var first time.Time
	var selected Report
	for _, report := range reports {
		if report.Status != "confirmed" || !strings.EqualFold(report.Country, "Israel") {
			continue
		}
		parsed, err := time.ParseInLocation("2006-01-02", report.Date, start.Location())
		if err != nil || parsed.Before(start) || parsed.After(end) {
			continue
		}
		if first.IsZero() || parsed.Before(first) {
			first, selected = parsed, report
		}
	}
	return first, selected, !first.IsZero()
}
