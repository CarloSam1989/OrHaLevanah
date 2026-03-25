package biblical

import (
	"time"

	"or-halevanah/internal/location"
)

func normalizeYearAnchor(year int, loc *time.Location) time.Time {
	return time.Date(year, time.March, 20, 12, 0, 0, 0, loc)
}

// findFirstMonthStartForGregorianYear estima el inicio del mes 1 (Aviv)
// buscando la primera cabeza de mes observacional desde una ventana
// que empieza cerca del equinoccio de primavera.
func findFirstMonthStartForGregorianYear(gregorianYear int, loc *time.Location) time.Time {
	anchor := normalizeYearAnchor(gregorianYear, loc)

	for i := 0; i <= 45; i++ {
		candidate := anchor.AddDate(0, 0, i)
		start := FindEstimatedMonthStart(candidate.Add(12 * time.Hour))
		if start.Equal(time.Date(candidate.Year(), candidate.Month(), candidate.Day(), 0, 0, 0, 0, loc)) {
			return start
		}
	}

	return time.Date(anchor.Year(), anchor.Month(), anchor.Day(), 0, 0, 0, 0, loc)
}

// EstimateBiblicalYearStart devuelve el inicio estimado del año bíblico vigente.
func EstimateBiblicalYearStart(now time.Time) time.Time {
	loc := now.Location()
	currentYearStart := findFirstMonthStartForGregorianYear(now.Year(), loc)

	if now.Before(currentYearStart) {
		return findFirstMonthStartForGregorianYear(now.Year()-1, loc)
	}

	return currentYearStart
}

// GetCurrentBiblicalMonthAt calcula el mes bíblico actual para una fecha dada.
func GetCurrentBiblicalMonthAt(now time.Time) int {
	yearStart := EstimateBiblicalYearStart(now)

	if now.Before(yearStart) {
		return 12
	}

	month := 1
	currentStart := yearStart

	for month < 12 {
		nextStart := FindNextEstimatedMonthStart(currentStart)
		if now.Before(nextStart) {
			return month
		}
		currentStart = nextStart
		month++
	}

	return 12
}

// Mantiene compatibilidad con el código existente.
func GetCurrentBiblicalMonth() int {
	jerusalem := location.GetJerusalem()

	loc, err := time.LoadLocation(jerusalem.Timezone)
	if err != nil {
		return 1
	}

	nowJerusalem := time.Now().In(loc)
	return GetCurrentBiblicalMonthAt(nowJerusalem)
}