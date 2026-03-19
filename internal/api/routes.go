package api

import "net/http"

func RegisterRoutes() {
	http.HandleFunc("/api/health", HealthHandler)
	http.HandleFunc("/api/jerusalem", JerusalemHandler)
}
