package biblical

import (
	"time"

	"or-halevanah/internal/lunar"
	"or-halevanah/internal/solar"
)

func normalizeDate(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
}

// isCandidateMonthStart evalúa si ese día civil, al atardecer en Jerusalén,
// puede ser cabeza de mes observacional.
// El día 1 bíblico comienza al atardecer de esa fecha civil.
func isCandidateMonthStart(day time.Time) bool {
	day = normalizeDate(day)
	loc := day.Location()

	sunset := solar.ApproxSunsetJerusalem(day)

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

// FindEstimatedMonthStart busca la cabeza de mes vigente.
// Si aún no cae el sol hoy, el mes vigente no puede comenzar hoy al atardecer,
// así que se busca hasta ayer.
func FindEstimatedMonthStart(now time.Time) time.Time {
	now = now.In(now.Location())
	today := normalizeDate(now)
	todaySunset := solar.ApproxSunsetJerusalem(today)

	searchEnd := today
	if now.Before(todaySunset) {
		searchEnd = today.AddDate(0, 0, -1)
	}

	// Buscar hacia atrás en una ventana amplia de 35 días
	// para no romper el cálculo del mes actual.
	for i := 0; i <= 35; i++ {
		candidate := searchEnd.AddDate(0, 0, -i)
		if isCandidateMonthStart(candidate) {
			return normalizeDate(candidate)
		}
	}

	// Fallback conservador
	return normalizeDate(searchEnd)
}

// FindNextEstimatedMonthStart busca el próximo posible inicio de mes
// a partir del mes actual. Se limita a la ventana normal de 29-30 días
// con pequeño margen.
func FindNextEstimatedMonthStart(currentMonthStart time.Time) time.Time {
	currentMonthStart = normalizeDate(currentMonthStart)

	// Un mes lunar observacional normal cae alrededor de 29 o 30 días.
	// Dejamos margen por seguridad.
	for i := 28; i <= 32; i++ {
		candidate := currentMonthStart.AddDate(0, 0, i)
		if isCandidateMonthStart(candidate) {
			return normalizeDate(candidate)
		}
	}

	// Fallback estable
	return normalizeDate(currentMonthStart.AddDate(0, 0, 30))
}

// CalculateBiblicalDay calcula el día bíblico actual.
// monthStart es la FECHA CIVIL cuyo atardecer inicia el día 1.
func CalculateBiblicalDay(monthStart, now time.Time, afterSunset bool) int {
	monthStart = normalizeDate(monthStart)
	monthStartSunset := solar.ApproxSunsetJerusalem(monthStart)

	// Si aún no ha comenzado el mes, protegemos el cálculo.
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