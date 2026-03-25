package biblical

import (
	"time"

	"or-halevanah/internal/lunar"
	"or-halevanah/internal/solar"
)

func normalizeDate(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
}

func isCandidateMonthStart(day time.Time) bool {
	day = normalizeDate(day)
	loc := day.Location()

	sunset := solar.ApproxSunsetJerusalem(day)

	// Revisamos poco después del atardecer.
	checkTime := time.Date(
		day.Year(),
		day.Month(),
		day.Day(),
		sunset.Hour(),
		sunset.Minute()+30,
		0,
		0,
		loc,
	)

	moon := lunar.CalculateMoonInfo(checkTime)
	ok, _ := lunar.IsPossibleNewMonthObservational(moon.AgeDays, true)

	return ok
}

// FindEstimatedMonthStart encuentra la cabeza de mes vigente.
// Si hoy todavía no llega el atardecer, el mes vigente no puede empezar hoy,
// por eso se busca hasta ayer.
func FindEstimatedMonthStart(now time.Time) time.Time {
	now = now.In(now.Location())
	today := normalizeDate(now)
	todaySunset := solar.ApproxSunsetJerusalem(today)

	searchEnd := today
	if now.Before(todaySunset) {
		searchEnd = today.AddDate(0, 0, -1)
	}

	for i := 0; i <= 35; i++ {
		candidate := searchEnd.AddDate(0, 0, -i)
		if isCandidateMonthStart(candidate) {
			return normalizeDate(candidate)
		}
	}

	// Fallback conservador
	return normalizeDate(searchEnd)
}

// FindNextEstimatedMonthStart busca la próxima cabeza de mes estimada
// partiendo del inicio actual.
func FindNextEstimatedMonthStart(currentMonthStart time.Time) time.Time {
	currentMonthStart = normalizeDate(currentMonthStart)

	// Ventana típica de mes lunar observacional: 29 o 30 días,
	// con un margen pequeño.
	for i := 29; i <= 31; i++ {
		candidate := currentMonthStart.AddDate(0, 0, i)
		if isCandidateMonthStart(candidate) {
			return normalizeDate(candidate)
		}
	}

	for i := 28; i <= 32; i++ {
		candidate := currentMonthStart.AddDate(0, 0, i)
		if isCandidateMonthStart(candidate) {
			return normalizeDate(candidate)
		}
	}

	return normalizeDate(currentMonthStart.AddDate(0, 0, 30))
}

// CalculateBiblicalDay calcula el día bíblico actual.
// monthStart es la FECHA CIVIL cuyo atardecer inicia el día 1.
func CalculateBiblicalDay(monthStart, now time.Time, afterSunset bool) int {
	monthStart = normalizeDate(monthStart)
	monthStartSunset := solar.ApproxSunsetJerusalem(monthStart)

	if now.Before(monthStartSunset) {
		return 1
	}

	diffHours := now.Sub(monthStartSunset).Hours()
	day := int(diffHours/24) + 1

	if day < 1 {
		return 1
	}
	if day > 30 {
		return 30
	}

	return day
}