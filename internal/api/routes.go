package api

import "net/http"

func RegisterRoutes() {
	http.HandleFunc("/api/health", HealthHandler)
	http.HandleFunc("/api/jerusalem", JerusalemHandler)
	http.HandleFunc("/api/moon/jerusalem/today", MoonJerusalemTodayHandler)
}
