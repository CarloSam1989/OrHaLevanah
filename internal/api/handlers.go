package api

import (
	"encoding/json"
	"net/http"
	"or-halevanah/internal/location"
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
