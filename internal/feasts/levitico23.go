package feasts

import (
	"time"

	"or-halevanah/internal/models"
)

func formatFeastDate(t time.Time) (string, string, string) {
	return t.Format("2006-01-02"), t.Weekday().String(), t.Format("2006-01-02 15:04")
}

func buildFixedFeast(name string, biblicalMonth, biblicalDay int, monthStart time.Time, description string) models.Feast {
	feastDate := monthStart.AddDate(0, 0, biblicalDay-1)
	gregorianDate, weekday, startAt := formatFeastDate(feastDate)

	return models.Feast{
		Name:            name,
		BiblicalMonth:   biblicalMonth,
		BiblicalDay:     biblicalDay,
		GregorianDate:   gregorianDate,
		Weekday:         weekday,
		BiblicalStartAt: startAt,
		Description:     description,
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
