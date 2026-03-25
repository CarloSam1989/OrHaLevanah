package biblical

import (
	"time"

	"or-halevanah/internal/location"
)

// findFirstMonthStartForGregorianYear estima el inicio del mes 1 (Aviv)
// buscando la primera cabeza de mes observacional a partir del 20 de marzo.
// Es una aproximación estable para no romper el sistema actual.
func findFirstMonthStartForGregorianYear(gregorianYear int, loc *time.Location) time.Time {
	anchor := time.Date(gregorianYear, time.March, 20, 12, 0, 0, 0, loc)

	// Buscar desde el 20 de marzo en adelante una ventana razonable.
	for i := 0; i <= 40; i++ {
		candidate := anchor.AddDate(0, 0, i)
		start := FindEstimatedMonthStart(candidate.Add(12 * time.Hour))
		if start.Equal(normalizeDate(candidate)) {
			return start
		}
	}

	// Fallback conservador
	return normalizeDate(anchor)
}

// EstimateBiblicalYearStart devuelve el inicio estimado del año bíblico
// vigente para la fecha dada.
func EstimateBiblicalYearStart(now time.Time) time.Time {
	loc := now.Location()

	currentYearStart := findFirstMonthStartForGregorianYear(now.Year(), loc)

	// Si aún no hemos llegado al inicio del año bíblico de este año gregoriano,
	// entonces seguimos en el año bíblico anterior.
	if now.Before(currentYearStart) {
		return findFirstMonthStartForGregorianYear(now.Year()-1, loc)
	}

	return currentYearStart
}

// GetCurrentBiblicalMonthAt calcula el mes bíblico actual para una fecha dada.
func GetCurrentBiblicalMonthAt(now time.Time) int {
	yearStart := EstimateBiblicalYearStart(now)

	// Si por alguna razón now cae antes del inicio calculado.
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

// GetCurrentBiblicalMonth mantiene la firma antigua para no romper handlers
// ni llamadas existentes.
func GetCurrentBiblicalMonth() int {
	jerusalem := location.GetJerusalem()

	loc, err := time.LoadLocation(jerusalem.Timezone)
	if err != nil {
		return 1
	}

	nowJerusalem := time.Now().In(loc)
	return GetCurrentBiblicalMonthAt(nowJerusalem)
}