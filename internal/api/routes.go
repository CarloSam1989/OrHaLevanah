package api

import "net/http"

func RegisterRoutes() {
	http.HandleFunc("/api/health", HealthHandler)
	http.HandleFunc("/api/jerusalem", JerusalemHandler)
	http.HandleFunc("/api/moon/jerusalem/today", MoonJerusalemTodayHandler)
	http.HandleFunc("/api/biblical/jerusalem/today", BiblicalJerusalemTodayHandler)
	http.HandleFunc("/api/biblical/jerusalem/month", BiblicalJerusalemMonthHandler)
	http.HandleFunc("/api/feasts/jerusalem", FeastsJerusalemHandler)

	http.HandleFunc("/api/aviv/status", AvivStatusHandler)
	http.HandleFunc("/api/aviv/reports", AvivReportsHandler)
	http.HandleFunc("/api/aviv/report", AvivReportHandler)
	http.HandleFunc("/api/admin/aviv/confirm", AvivConfirmHandler)
	http.HandleFunc("/api/admin/aviv/reject", AvivRejectHandler)
	http.HandleFunc("/api/admin/aviv/reset", AvivResetHandler)
}
