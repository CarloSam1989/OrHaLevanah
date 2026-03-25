package solar

import (
	"math"
	"time"
)

// ApproxSunsetJerusalem devuelve una aproximación diaria del atardecer en Jerusalén.
// No es un cálculo astronómico exacto, pero es más estable que usar una sola hora fija por mes.
// Rango aproximado anual usado:
// - invierno: ~16:36
// - verano:   ~18:49
func ApproxSunsetJerusalem(t time.Time) time.Time {
	year, month, day := t.Date()
	loc := t.Location()

	dayOfYear := t.YearDay()

	// Onda anual simple.
	// Ajuste para que el mínimo caiga cerca de diciembre-enero y el máximo cerca de junio-julio.
	angle := 2 * math.Pi * float64(dayOfYear-172) / 365.2422

	// Centro anual aproximado: 17:42.5 => 1062.5 minutos
	// Amplitud aproximada: 66.5 minutos
	// Así obtenemos aprox entre 16:36 y 18:49
	minutes := 1062.5 + 66.5*math.Cos(angle)

	hour := int(minutes) / 60
	minute := int(math.Round(minutes)) % 60

	if minute >= 60 {
		minute -= 60
		hour++
	}
	if minute < 0 {
		minute += 60
		hour--
	}

	return time.Date(year, month, day, hour, minute, 0, 0, loc)
}