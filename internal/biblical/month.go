package biblical

import (
	"time"

	"or-halevanah/internal/lunar"
	"or-halevanah/internal/solar"
)

// FindEstimatedMonthStart busca hacia atrás el inicio estimado del mes bíblico.
func FindEstimatedMonthStart(now time.Time) time.Time {
	loc := now.Location()

	for i := 0; i < 5; i++ {
		day := now.AddDate(0, 0, -i)

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

		if ok {
			return time.Date(day.Year(), day.Month(), day.Day(), 0, 0, 0, 0, loc)
		}
	}

	return time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, loc)
}

// CalculateBiblicalDay calcula el día del mes bíblico actual.
func CalculateBiblicalDay(monthStart, now time.Time, afterSunset bool) int {
	startDate := time.Date(monthStart.Year(), monthStart.Month(), monthStart.Day(), 0, 0, 0, 0, monthStart.Location())
	currentDate := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())

	diffDays := int(currentDate.Sub(startDate).Hours() / 24)

	day := diffDays + 1
	if afterSunset {
		day++
	}

	if day < 1 {
		return 1
	}
	if day > 30 {
		return 30
	}

	return day
}
