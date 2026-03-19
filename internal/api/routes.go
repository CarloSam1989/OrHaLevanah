package api

import "net/http"

func RegisterRoutes() {
	http.HandleFunc("/api/health", HealthHandler)
	http.HandleFunc("/api/jerusalem", JerusalemHandler)
	http.HandleFunc("/api/moon/jerusalem/today", MoonJerusalemTodayHandler)
	http.HandleFunc("/api/biblical/jerusalem/today", BiblicalJerusalemTodayHandler)
	http.HandleFunc("/api/biblical/jerusalem/month", BiblicalJerusalemMonthHandler)
}
