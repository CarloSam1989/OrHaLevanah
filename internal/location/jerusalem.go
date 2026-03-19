package location

import "or-halevanah/internal/models"

func GetJerusalem() models.Location {
	return models.Location{
		Name:      "Jerusalem",
		Latitude:  31.7683,
		Longitude: 35.2137,
		Timezone:  "Asia/Jerusalem",
	}
}
