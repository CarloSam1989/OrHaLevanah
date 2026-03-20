package feasts

import (
	"time"

	"or-halevanah/internal/models"
)

func FindFirstShabbatInUnleavenedBread(monthStart time.Time) time.Time {
	for day := 15; day <= 21; day++ {
		current := monthStart.AddDate(0, 0, day-1)
		if current.Weekday() == time.Saturday {
			return current
		}
	}
	return time.Time{}
}

func CalculateBikkurim(monthStart time.Time) models.Feast {
	shabbat := FindFirstShabbatInUnleavenedBread(monthStart)
	bikkurim := shabbat.AddDate(0, 0, 1)

	return models.Feast{
		Name:            "Bikkurim / Primeros Frutos",
		BiblicalMonth:   1,
		BiblicalDay:     int(bikkurim.Sub(monthStart).Hours()/24) + 1,
		GregorianDate:   bikkurim.Format("2006-01-02"),
		Weekday:         bikkurim.Weekday().String(),
		BiblicalStartAt: bikkurim.Format("2006-01-02 15:04"),
		Description:     "Día después del primer Shabat semanal dentro de Panes sin Levadura. Aquí inicia el conteo del Omer. El día bíblico comienza al atardecer.",
	}
}

func CalculateShavuot(monthStart time.Time) models.Feast {
	bikkurim := CalculateBikkurim(monthStart)
	bikkurimDate, _ := time.Parse("2006-01-02", bikkurim.GregorianDate)
	shavuot := bikkurimDate.AddDate(0, 0, 49)

	// Shavuot puede caer ya en otro mes bíblico; aquí dejamos el mes/día bíblico como referencia aproximada anual.
	return models.Feast{
		Name:            "Shavuot",
		BiblicalMonth:   0,
		BiblicalDay:     0,
		GregorianDate:   shavuot.Format("2006-01-02"),
		Weekday:         shavuot.Weekday().String(),
		BiblicalStartAt: shavuot.Format("2006-01-02 15:04"),
		Description:     "Shavuot cae 50 días contando desde Bikkurim. El día bíblico comienza al atardecer.",
	}
}

func CalculateOmerInfo(monthStart, nowJerusalem time.Time) models.OmerInfo {
	shabbat := FindFirstShabbatInUnleavenedBread(monthStart)
	bikkurim := shabbat.AddDate(0, 0, 1)
	shavuot := bikkurim.AddDate(0, 0, 49)

	currentCivilDay := time.Date(
		nowJerusalem.Year(),
		nowJerusalem.Month(),
		nowJerusalem.Day(),
		0, 0, 0, 0,
		nowJerusalem.Location(),
	)

	omerDay := 0
	isCounting := false

	if !currentCivilDay.Before(bikkurim) && !currentCivilDay.After(shavuot) {
		omerDay = int(currentCivilDay.Sub(bikkurim).Hours()/24) + 1
		if omerDay >= 1 && omerDay <= 50 {
			isCounting = true
		}
	}

	return models.OmerInfo{
		FirstShabbatInMatzot: shabbat.Format("2006-01-02"),
		BikkurimDate:         bikkurim.Format("2006-01-02"),
		BikkurimWeekday:      bikkurim.Weekday().String(),
		ShavuotDate:          shavuot.Format("2006-01-02"),
		ShavuotWeekday:       shavuot.Weekday().String(),
		OmerDayToday:         omerDay,
		IsOmerCounting:       isCounting,
		DayNote:              "Las fechas bíblicas comienzan al atardecer y continúan hasta el siguiente atardecer.",
	}
}

func GetCurrentFeasts(monthStart, nowJerusalem time.Time) []models.Feast {
	var result []models.Feast

	today := nowJerusalem.Format("2006-01-02")

	for _, feast := range GetMonthOneFixedFeasts(monthStart) {
		if feast.GregorianDate == today {
			result = append(result, feast)
		}
	}

	bikkurim := CalculateBikkurim(monthStart)
	if bikkurim.GregorianDate == today {
		result = append(result, bikkurim)
	}

	shavuot := CalculateShavuot(monthStart)
	if shavuot.GregorianDate == today {
		result = append(result, shavuot)
	}

	return result
}

func GetUpcomingFeasts(monthStart, nowJerusalem time.Time) []models.Feast {
	var result []models.Feast

	today := nowJerusalem.Format("2006-01-02")

	for _, feast := range GetMonthOneFixedFeasts(monthStart) {
		if feast.GregorianDate > today {
			result = append(result, feast)
		}
	}

	bikkurim := CalculateBikkurim(monthStart)
	if bikkurim.GregorianDate > today {
		result = append(result, bikkurim)
	}

	shavuot := CalculateShavuot(monthStart)
	if shavuot.GregorianDate > today {
		result = append(result, shavuot)
	}

	return result
}
