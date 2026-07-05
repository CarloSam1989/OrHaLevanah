package feasts

import (
	"testing"
	"time"

	"or-halevanah/internal/models"
)

func TestCalculateOmerInfoHasOneAnnualWindow(t *testing.T) {
	loc, err := time.LoadLocation("Asia/Jerusalem")
	if err != nil {
		t.Fatal(err)
	}

	avivStart := time.Date(2026, time.March, 20, 0, 0, 0, 0, loc)
	omer := CalculateOmerInfo(avivStart, avivStart)
	bikkurim, err := time.ParseInLocation("2006-01-02", omer.BikkurimDate, loc)
	if err != nil {
		t.Fatal(err)
	}

	during := CalculateOmerInfo(avivStart, bikkurim.Add(20*time.Hour))
	if !during.IsOmerCounting || during.OmerDayToday != 1 {
		t.Fatalf("se esperaba el día 1 activo, se obtuvo activo=%v día=%d", during.IsOmerCounting, during.OmerDayToday)
	}

	after := CalculateOmerInfo(avivStart, bikkurim.AddDate(0, 0, 60))
	if after.IsOmerCounting || after.OmerDayToday != 0 {
		t.Fatalf("el conteo no debe repetirse fuera de su ventana anual: activo=%v día=%d", after.IsOmerCounting, after.OmerDayToday)
	}
}

func TestAllFeastsIncludesSeventhMonthVisualEvents(t *testing.T) {
	loc, err := time.LoadLocation("Asia/Jerusalem")
	if err != nil {
		t.Fatal(err)
	}

	avivStart := time.Date(2026, time.March, 20, 0, 0, 0, 0, loc)
	all := GetAllFeasts(avivStart)
	found := map[string]models.Feast{}

	for _, feast := range all {
		if feast.Name == "Yom Teruah" || feast.Name == "Yom Kippur" || feast.Name == "Sucot" {
			found[feast.Name] = feast
		}
	}

	for _, name := range []string{"Yom Teruah", "Yom Kippur", "Sucot"} {
		if _, ok := found[name]; !ok {
			t.Fatalf("falta %s en las fiestas anuales", name)
		}
	}

	start, err := time.Parse("2006-01-02", found["Sucot"].GregorianStartDate)
	if err != nil {
		t.Fatal(err)
	}
	end, err := time.Parse("2006-01-02", found["Sucot"].GregorianEndDate)
	if err != nil {
		t.Fatal(err)
	}
	if days := int(end.Sub(start).Hours() / 24); days != 7 {
		t.Fatalf("Sucot debe abarcar siete días completos; obtuvo %d", days)
	}
}
