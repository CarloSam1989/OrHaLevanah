package horizons

import (
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"or-halevanah/internal/models"
	"or-halevanah/internal/solar"
)

const (
	defaultAPIURL    = "https://ssd.jpl.nasa.gov/api/horizons.api"
	defaultCachePath = "data/horizons_jerusalem.json"
)

type cacheFile struct {
	Forecast models.HorizonsForecast `json:"forecast"`
}

type horizonsResponse struct {
	Result    string `json:"result"`
	Error     string `json:"error"`
	Signature struct {
		Source  string `json:"source"`
		Version string `json:"version"`
	} `json:"signature"`
}

var (
	mu          sync.RWMutex
	lastAttempt time.Time
	refreshTTL  = 6 * time.Hour
)

func apiURL() string {
	if value := strings.TrimSpace(os.Getenv("JPL_HORIZONS_API_URL")); value != "" {
		return value
	}
	return defaultAPIURL
}

func cachePath() string {
	if value := strings.TrimSpace(os.Getenv("JPL_HORIZONS_CACHE_PATH")); value != "" {
		return value
	}
	return defaultCachePath
}

func LoadForecast() (models.HorizonsForecast, error) {
	mu.RLock()
	defer mu.RUnlock()
	content, err := os.ReadFile(cachePath())
	if err != nil {
		return models.HorizonsForecast{}, err
	}
	var cache cacheFile
	if err := json.Unmarshal(content, &cache); err != nil {
		return models.HorizonsForecast{}, err
	}
	return cache.Forecast, nil
}

func ForecastFor(currentMonthStart time.Time) (models.HorizonsForecast, bool) {
	forecast, err := LoadForecast()
	if err != nil || forecast.BasedOnMonthStart != currentMonthStart.Format("2006-01-02") {
		return models.HorizonsForecast{}, false
	}
	return forecast, forecast.PossibleRoshChodesh != ""
}

func RefreshIfStale(currentMonthStart time.Time, location models.Location) error {
	mu.Lock()
	if time.Since(lastAttempt) < refreshTTL {
		mu.Unlock()
		return nil
	}
	lastAttempt = time.Now()
	mu.Unlock()

	if cached, ok := ForecastFor(currentMonthStart); ok && time.Since(cached.CheckedAt) < refreshTTL {
		return nil
	}

	forecast, err := requestForecast(currentMonthStart, location)
	if err != nil {
		return err
	}
	return saveForecast(forecast)
}

func requestForecast(currentMonthStart time.Time, location models.Location) (models.HorizonsForecast, error) {
	loc, err := time.LoadLocation(location.Timezone)
	if err != nil {
		return models.HorizonsForecast{}, err
	}
	times := make([]string, 0, 35)
	for days := 28; days <= 32; days++ {
		day := currentMonthStart.In(loc).AddDate(0, 0, days)
		sunset := solar.ApproxSunsetJerusalem(day)
		// Sampleamos alrededor del ocaso aproximado. Horizons indica con "C"
		// cuales filas estan realmente dentro del crepusculo civil.
		for offset := -20; offset <= 40; offset += 10 {
			sample := sunset.Add(time.Duration(offset) * time.Minute).UTC()
			times = append(times, "'"+sample.Format("2006-Jan-02 15:04")+"'")
		}
	}

	params := url.Values{}
	params.Set("format", "json")
	params.Set("COMMAND", "'301'")
	params.Set("OBJ_DATA", "'NO'")
	params.Set("MAKE_EPHEM", "'YES'")
	params.Set("EPHEM_TYPE", "'OBSERVER'")
	params.Set("CENTER", "'coord@399'")
	params.Set("COORD_TYPE", "'GEODETIC'")
	params.Set("SITE_COORD", fmt.Sprintf("'%.4f,%.4f,0.754'", location.Longitude, location.Latitude))
	params.Set("TLIST", strings.Join(times, ","))
	params.Set("TLIST_TYPE", "'CAL'")
	params.Set("TIME_TYPE", "'UT'")
	params.Set("QUANTITIES", "'4,10,23,24'")
	params.Set("CSV_FORMAT", "'YES'")
	params.Set("ANG_FORMAT", "'DEG'")
	params.Set("APPARENT", "'REFRACTED'")

	req, err := http.NewRequest(http.MethodGet, apiURL()+"?"+params.Encode(), nil)
	if err != nil {
		return models.HorizonsForecast{}, err
	}
	req.Header.Set("User-Agent", "Or-HaLevanah/1.0 JPL-Horizons-client")
	response, err := (&http.Client{Timeout: 20 * time.Second}).Do(req)
	if err != nil {
		return models.HorizonsForecast{}, err
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		return models.HorizonsForecast{}, fmt.Errorf("JPL Horizons returned HTTP %d", response.StatusCode)
	}
	body, err := io.ReadAll(io.LimitReader(response.Body, 2<<20))
	if err != nil {
		return models.HorizonsForecast{}, err
	}
	var payload horizonsResponse
	if err := json.Unmarshal(body, &payload); err != nil {
		return models.HorizonsForecast{}, err
	}
	if payload.Error != "" {
		return models.HorizonsForecast{}, errors.New(payload.Error)
	}
	forecast, err := parseEphemeris(payload.Result, loc)
	if err != nil {
		return models.HorizonsForecast{}, err
	}
	forecast.Source = "NASA/JPL Horizons"
	forecast.SourceURL = "https://ssd.jpl.nasa.gov/horizons/"
	forecast.APIVersion = payload.Signature.Version
	forecast.CheckedAt = time.Now().UTC()
	forecast.BasedOnMonthStart = currentMonthStart.Format("2006-01-02")
	return forecast, nil
}

func parseEphemeris(result string, loc *time.Location) (models.HorizonsForecast, error) {
	start := strings.Index(result, "$$SOE")
	end := strings.Index(result, "$$EOE")
	if start < 0 || end <= start {
		return models.HorizonsForecast{}, errors.New("JPL Horizons response has no ephemeris rows")
	}
	reader := csv.NewReader(strings.NewReader(result[start+len("$$SOE") : end]))
	reader.TrimLeadingSpace = true
	reader.FieldsPerRecord = -1
	for {
		row, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil || len(row) < 7 {
			continue
		}
		for i := range row {
			row[i] = strings.TrimSpace(row[i])
		}
		for len(row) > 0 && row[len(row)-1] == "" {
			row = row[:len(row)-1]
		}
		if len(row) < 7 {
			continue
		}
		// The requested quantities end in AZ, EL, illumination, elongation,
		// trailing/leading marker and phase angle. Reading from the end also
		// tolerates the separate solar/lunar presence columns at the beginning.
		phase, e1 := strconv.ParseFloat(row[len(row)-1], 64)
		elongation, e2 := strconv.ParseFloat(row[len(row)-3], 64)
		illumination, e3 := strconv.ParseFloat(row[len(row)-4], 64)
		elevation, e4 := strconv.ParseFloat(row[len(row)-5], 64)
		azimuth, e5 := strconv.ParseFloat(row[len(row)-6], 64)
		if e1 != nil || e2 != nil || e3 != nil || e4 != nil || e5 != nil {
			continue
		}
		marker := row[len(row)-2]
		observedAt, err := parseHorizonsTime(row[0], loc)
		if err != nil {
			continue
		}
		// Conservative first-crescent screen. Horizons provides geometry,
		// while weather and a human report still determine confirmation.
		solarPresence := ""
		if len(row) > 1 {
			solarPresence = row[1]
		}
		if solarPresence == "C" && strings.Contains(marker, "/T") && elevation >= 3 && elongation >= 10 && illumination >= 0.5 {
			return models.HorizonsForecast{
				PossibleRoshChodesh: observedAt.Format("2006-01-02"),
				SampleTimeJerusalem: observedAt.Format("2006-01-02 15:04"),
				AzimuthDegrees:      azimuth,
				ElevationDegrees:    elevation,
				IlluminationPercent: illumination,
				SolarElongation:     elongation,
				PhaseAngle:          phase,
				EveningCrescent:     true,
				Status:              "estimated",
				Note:                "Posible creciente visible desde Jerusalén; pendiente de confirmación humana en Israel.",
			}, nil
		}
	}
	return models.HorizonsForecast{}, errors.New("no conservative crescent window found in JPL Horizons response")
}

func parseHorizonsTime(value string, loc *time.Location) (time.Time, error) {
	value = strings.TrimSpace(strings.TrimPrefix(value, "A.D."))
	for _, layout := range []string{"2006-Jan-02 15:04", "2006-Jan-02 15:04:05"} {
		parsed, err := time.ParseInLocation(layout, value, time.UTC)
		if err == nil {
			return parsed.In(loc), nil
		}
	}
	return time.Time{}, fmt.Errorf("invalid Horizons time %q", value)
}

func saveForecast(forecast models.HorizonsForecast) error {
	mu.Lock()
	defer mu.Unlock()
	content, err := json.MarshalIndent(cacheFile{Forecast: forecast}, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(cachePath(), append(content, '\n'), 0644)
}
