package api

import (
	"encoding/json"
	"net/http"
	"time"

	"or-halevanah/internal/location"
	"or-halevanah/internal/lunar"
	"or-halevanah/internal/models"
)

func HealthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	response := models.ApiResponse{
		Success: true,
		Message: "API activa",
		Data: map[string]string{
			"service": "Or-HaLevanah",
			"status":  "ok",
		},
	}

	json.NewEncoder(w).Encode(response)
}

func JerusalemHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	jerusalem := location.GetJerusalem()

	response := models.ApiResponse{
		Success: true,
		Message: "Ubicación de referencia cargada correctamente",
		Data:    jerusalem,
	}

	json.NewEncoder(w).Encode(response)
}

func MoonJerusalemTodayHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	jerusalem := location.GetJerusalem()

	loc, err := time.LoadLocation(jerusalem.Timezone)
	if err != nil {
		http.Error(w, "No se pudo cargar la zona horaria de Jerusalén", http.StatusInternalServerError)
		return
	}

	nowJerusalem := time.Now().In(loc)
	moonInfo := lunar.CalculateMoonInfo(nowJerusalem)

	data := models.MoonTodayResponse{
		Date:     nowJerusalem.Format("2006-01-02"),
		Location: jerusalem,
		Moon:     moonInfo,
	}

	response := models.ApiResponse{
		Success: true,
		Message: "Información lunar calculada correctamente",
		Data:    data,
	}

	json.NewEncoder(w).Encode(response)
}
