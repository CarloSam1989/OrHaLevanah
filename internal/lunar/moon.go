package lunar

import (
	"math"
	"time"

	"or-halevanah/internal/models"
)

const synodicMonth = 29.53058867

var knownNewMoon = time.Date(2000, 1, 6, 18, 14, 0, 0, time.UTC)

func CalculateMoonInfo(t time.Time) models.MoonInfo {
	t = t.UTC()

	age := moonAgeDays(t)
	illumination := moonIllumination(age)
	phase := moonPhaseName(age)
	lunarDay := int(math.Floor(age)) + 1

	if lunarDay < 1 {
		lunarDay = 1
	}
	if lunarDay > 30 {
		lunarDay = 30
	}

	return models.MoonInfo{
		AgeDays:      round(age, 2),
		Illumination: round(illumination, 4),
		PhaseName:    phase,
		LunarDay:     lunarDay,
	}
}

func IsPossibleNewMonthObservational(ageDays float64, afterSunset bool) (bool, string) {
	if !afterSunset {
		return false, "Aún no es después del atardecer en Jerusalén"
	}

	if ageDays >= 1.0 && ageDays <= 2.5 {
		return true, "Posible inicio observacional del mes: luna joven visible potencial después del atardecer"
	}

	if ageDays < 1.0 {
		return false, "La luna es demasiado joven para una observación inicial confiable"
	}

	return false, "La ventana principal de observación inicial del creciente probablemente ya pasó"
}

func moonAgeDays(t time.Time) float64 {
	diff := t.Sub(knownNewMoon).Hours() / 24.0
	age := math.Mod(diff, synodicMonth)
	if age < 0 {
		age += synodicMonth
	}
	return age
}

func moonIllumination(age float64) float64 {
	phaseAngle := 2 * math.Pi * age / synodicMonth
	return (1 - math.Cos(phaseAngle)) / 2
}

func moonPhaseName(age float64) string {
	switch {
	case age < 1.84566:
		return "New Moon"
	case age < 5.53699:
		return "Waxing Crescent"
	case age < 9.22831:
		return "First Quarter"
	case age < 12.91963:
		return "Waxing Gibbous"
	case age < 16.61096:
		return "Full Moon"
	case age < 20.30228:
		return "Waning Gibbous"
	case age < 23.99361:
		return "Last Quarter"
	case age < 27.68493:
		return "Waning Crescent"
	default:
		return "New Moon"
	}
}

func round(value float64, places int) float64 {
	factor := math.Pow(10, float64(places))
	return math.Round(value*factor) / factor
}
