package feasts

import (
	"time"

	"or-halevanah/internal/models"
	"or-halevanah/internal/solar"
)

func buildBiblicalWindow(day time.Time) (time.Time, time.Time) {
	startAt := solar.ApproxSunsetJerusalem(day)
	endAt := solar.ApproxSunsetJerusalem(day.AddDate(0, 0, 1))
	return startAt, endAt
}

func formatFeastDate(start, end time.Time) (string, string, string, string, string) {
	return start.Format("2006-01-02"),
		start.Weekday().String(),
		start.Format("2006-01-02 15:04"),
		end.Format("2006-01-02 15:04"),
		end.Format("2006-01-02")
}

func buildFixedFeast(name string, biblicalMonth, biblicalDay int, monthStart time.Time, description string) models.Feast {
	feastDay := monthStart.AddDate(0, 0, biblicalDay-1)

	startAt, endAt := buildBiblicalWindow(feastDay)

	return models.Feast{
		Name:               name,
		BiblicalMonth:      biblicalMonth,
		BiblicalDay:        biblicalDay,
		GregorianDate:      startAt.Format("2006-01-02"),
		Weekday:            startAt.Weekday().String(),
		BiblicalStartAt:    startAt.Format("2006-01-02 15:04"),
		BiblicalEndAt:      endAt.Format("2006-01-02 15:04"),
		GregorianStartDate: startAt.Format("2006-01-02"),
		GregorianEndDate:   endAt.Format("2006-01-02"),
		Description:        description,
	}
}

func GetMonthOneFixedFeasts(monthStart time.Time) []models.Feast {
	return []models.Feast{
		buildFixedFeast(
			"Pesaj",
			1,
			14,
			monthStart,
			"Pesaj cae en el día 14 del mes 1. El día bíblico comienza al atardecer.",
		),
		buildFixedFeast(
			"Panes sin Levadura - Inicio",
			1,
			15,
			monthStart,
			"Inicio de Panes sin Levadura. El día bíblico comienza al atardecer.",
		),
		buildFixedFeast(
			"Panes sin Levadura - Fin",
			1,
			21,
			monthStart,
			"Fin de Panes sin Levadura. El día bíblico comienza al atardecer.",
		),
	}
}

func GetSeventhMonthFixedFeasts(monthSevenStart time.Time) []models.Feast {
	return []models.Feast{
		buildFixedFeast(
			"Yom Teruah",
			7,
			1,
			monthSevenStart,
			"Fiesta de Yom Teruah. El día bíblico comienza al atardecer.",
		),
		buildFixedFeast(
			"Yom Kippur",
			7,
			10,
			monthSevenStart,
			"Día de expiación. El día bíblico comienza al atardecer.",
		),
		buildFixedFeast(
			"Sucot - Inicio",
			7,
			15,
			monthSevenStart,
			"Inicio de Sucot. El día bíblico comienza al atardecer.",
		),
		buildFixedFeast(
			"Sucot - Fin",
			7,
			21,
			monthSevenStart,
			"Fin de Sucot. El día bíblico comienza al atardecer.",
		),
		buildFixedFeast(
			"Shemini Atzeret",
			7,
			22,
			monthSevenStart,
			"Octavo día de asamblea solemne. El día bíblico comienza al atardecer.",
		),
	}
}
