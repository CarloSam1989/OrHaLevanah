package feasts

import (
	"sort"
	"time"

	"or-halevanah/internal/models"
	"or-halevanah/internal/solar"
)

func parseDate(dateStr string) time.Time {
	t, _ := time.Parse("2006-01-02", dateStr)
	return t
}

func parseDateTime(dateTimeStr string) time.Time {
	t, _ := time.Parse("2006-01-02 15:04", dateTimeStr)
	return t
}

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
	bikkurimDay := shabbat.AddDate(0, 0, 1)

	startAt := solar.ApproxSunsetJerusalem(bikkurimDay)
	endAt := solar.ApproxSunsetJerusalem(bikkurimDay.AddDate(0, 0, 1))

	return models.Feast{
		Name:               "Bikkurim / Primeros Frutos",
		BiblicalMonth:      1,
		BiblicalDay:        int(bikkurimDay.Sub(monthStart).Hours()/24) + 1,
		GregorianDate:      startAt.Format("2006-01-02"),
		Weekday:            startAt.Weekday().String(),
		BiblicalStartAt:    startAt.Format("2006-01-02 15:04"),
		BiblicalEndAt:      endAt.Format("2006-01-02 15:04"),
		GregorianStartDate: startAt.Format("2006-01-02"),
		GregorianEndDate:   endAt.Format("2006-01-02"),
		Description:        "Día después del primer Shabat semanal dentro de Panes sin Levadura. Aquí inicia el conteo del Omer. El día bíblico comienza al atardecer.",
	}
}

func CalculateShavuot(monthStart time.Time) models.Feast {
	bikkurim := CalculateBikkurim(monthStart)
	bikkurimDate, _ := time.Parse("2006-01-02", bikkurim.GregorianDate)
	shavuotDay := bikkurimDate.AddDate(0, 0, 49)

	startAt := solar.ApproxSunsetJerusalem(shavuotDay)
	endAt := solar.ApproxSunsetJerusalem(shavuotDay.AddDate(0, 0, 1))

	return models.Feast{
		Name:               "Shavuot",
		BiblicalMonth:      3,
		BiblicalDay:        6,
		GregorianDate:      startAt.Format("2006-01-02"),
		Weekday:            startAt.Weekday().String(),
		BiblicalStartAt:    startAt.Format("2006-01-02 15:04"),
		BiblicalEndAt:      endAt.Format("2006-01-02 15:04"),
		GregorianStartDate: startAt.Format("2006-01-02"),
		GregorianEndDate:   endAt.Format("2006-01-02"),
		Description:        "Shavuot ocurre después de completar el conteo de 7 semanas (49 días) desde Bikkurim. Es el día 50. El día bíblico comienza al atardecer.",
	}
}

func CalculateOmerInfo(monthStart, nowJerusalem time.Time) models.OmerInfo {
	shabbat := FindFirstShabbatInUnleavenedBread(monthStart)
	bikkurimDay := shabbat.AddDate(0, 0, 1)
	shavuotDay := bikkurimDay.AddDate(0, 0, 49)

	bikkurimStart := solar.ApproxSunsetJerusalem(bikkurimDay)
	shavuotEnd := solar.ApproxSunsetJerusalem(shavuotDay.AddDate(0, 0, 1))

	omerDay := 0
	isCounting := false

	if (nowJerusalem.Equal(bikkurimStart) || nowJerusalem.After(bikkurimStart)) && nowJerusalem.Before(shavuotEnd) {
		daysSinceStart := int(nowJerusalem.Sub(bikkurimStart).Hours() / 24)
		omerDay = daysSinceStart + 1

		if omerDay >= 1 && omerDay <= 50 {
			isCounting = true
		}
	}

	return models.OmerInfo{
		FirstShabbatInMatzot: shabbat.Format("2006-01-02"),
		BikkurimDate:         bikkurimDay.Format("2006-01-02"),
		BikkurimWeekday:      bikkurimDay.Weekday().String(),
		ShavuotDate:          shavuotDay.Format("2006-01-02"),
		ShavuotWeekday:       shavuotDay.Weekday().String(),
		OmerDayToday:         omerDay,
		IsOmerCounting:       isCounting,
		DayNote:              "Las fechas bíblicas comienzan al atardecer y continúan hasta el siguiente atardecer.",
	}
}

func GetAllFeasts(monthStart time.Time) []models.Feast {
	var all []models.Feast

	// Mes 1
	all = append(all, GetMonthOneFixedFeasts(monthStart)...)

	// Bikurim y Shavuot
	all = append(all, CalculateBikkurim(monthStart))
	all = append(all, CalculateShavuot(monthStart))

	// Mes 7
	monthSevenStart := monthStart.AddDate(0, 6, 0)
	all = append(all, GetSeventhMonthFixedFeasts(monthSevenStart)...)

	sort.Slice(all, func(i, j int) bool {
		return parseDate(all[i].GregorianDate).Before(parseDate(all[j].GregorianDate))
	})

	return all
}

func GetCurrentFeasts(monthStart, nowJerusalem time.Time) []models.Feast {
	result := make([]models.Feast, 0)

	for _, feast := range GetAllFeasts(monthStart) {
		startAt := parseDateTime(feast.BiblicalStartAt)
		endAt := parseDateTime(feast.BiblicalEndAt)

		if (nowJerusalem.Equal(startAt) || nowJerusalem.After(startAt)) && nowJerusalem.Before(endAt) {
			result = append(result, feast)
		}
	}

	return result
}

func GetUpcomingFeasts(monthStart, nowJerusalem time.Time) []models.Feast {
	result := make([]models.Feast, 0)

	for _, feast := range GetAllFeasts(monthStart) {
		startAt := parseDateTime(feast.BiblicalStartAt)

		if startAt.After(nowJerusalem) {
			result = append(result, feast)
		}
	}

	sort.Slice(result, func(i, j int) bool {
		return parseDateTime(result[i].BiblicalStartAt).Before(parseDateTime(result[j].BiblicalStartAt))
	})

	return result
}
