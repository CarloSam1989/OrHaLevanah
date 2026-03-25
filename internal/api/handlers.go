package api

import (
	"encoding/json"
	"net/http"
	"time"

	"or-halevanah/internal/aviv"
	"or-halevanah/internal/biblical"
	"or-halevanah/internal/feasts"
	"or-halevanah/internal/location"
	"or-halevanah/internal/lunar"
	"or-halevanah/internal/models"
	"or-halevanah/internal/solar"
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

func BiblicalJerusalemTodayHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	jerusalem := location.GetJerusalem()

	loc, err := time.LoadLocation(jerusalem.Timezone)
	if err != nil {
		http.Error(w, "No se pudo cargar la zona horaria de Jerusalén", http.StatusInternalServerError)
		return
	}

	nowJerusalem := time.Now().In(loc)
	sunset := solar.ApproxSunsetJerusalem(nowJerusalem)
	afterSunset := biblical.IsAfterSunset(nowJerusalem, sunset)
	biblicalDate := biblical.GetBiblicalDate(nowJerusalem, sunset)
	moonInfo := lunar.CalculateMoonInfo(nowJerusalem)

	possibleNewMonth, note := lunar.IsPossibleNewMonthObservational(moonInfo.AgeDays, afterSunset)

	data := models.BiblicalTodayResponse{
		CivilDate:     nowJerusalem.Format("2006-01-02"),
		BiblicalDate:  biblicalDate.Format("2006-01-02"),
		JerusalemTime: nowJerusalem.Format("2006-01-02 15:04:05"),
		SunsetTime:    sunset.Format("15:04"),
		AfterSunset:   afterSunset,
		DayNote:       "El día bíblico comienza al atardecer en Jerusalén y termina al siguiente atardecer.",
		Location:      jerusalem,
		Moon:          moonInfo,
		Observational: models.ObservationalInfo{
			PossibleNewMonth: possibleNewMonth,
			VisibilityNote:   note,
		},
	}

	response := models.ApiResponse{
		Success: true,
		Message: "Información bíblica de Jerusalén calculada correctamente",
		Data:    data,
	}

	json.NewEncoder(w).Encode(response)
}

func BiblicalJerusalemMonthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	jerusalem := location.GetJerusalem()

	loc, err := time.LoadLocation(jerusalem.Timezone)
	if err != nil {
		http.Error(w, "No se pudo cargar la zona horaria de Jerusalén", http.StatusInternalServerError)
		return
	}

	nowJerusalem := time.Now().In(loc)
	sunset := solar.ApproxSunsetJerusalem(nowJerusalem)
	afterSunset := biblical.IsAfterSunset(nowJerusalem, sunset)
	biblicalDate := biblical.GetBiblicalDate(nowJerusalem, sunset)

	moonInfo := lunar.CalculateMoonInfo(nowJerusalem)
	monthStart := biblical.FindEstimatedMonthStart(nowJerusalem)
	nextMonthStart := biblical.FindNextEstimatedMonthStart(monthStart)
	biblicalDay := biblical.CalculateBiblicalDay(monthStart, nowJerusalem, afterSunset)

	data := models.BiblicalMonthResponse{
		CivilDate:            nowJerusalem.Format("2006-01-02"),
		BiblicalDate:         biblicalDate.Format("2006-01-02"),
		JerusalemTime:        nowJerusalem.Format("2006-01-02 15:04:05"),
		SunsetTime:           sunset.Format("15:04"),
		AfterSunset:          afterSunset,
		DayNote:              "El día bíblico comienza al atardecer en Jerusalén y termina al siguiente atardecer.",
		MonthStart:           monthStart.Format("2006-01-02"),
		NextMonthStart:       nextMonthStart.Format("2006-01-02"),
		BiblicalDay:          biblicalDay,
		IsPossibleDayOne:     biblicalDay == 1,
		IsPossibleNextDayOne: false,
		Location:             jerusalem,
		Moon:                 moonInfo,
	}

	response := models.ApiResponse{
		Success: true,
		Message: "Mes bíblico calculado correctamente",
		Data:    data,
	}

	json.NewEncoder(w).Encode(response)
}

func FeastsJerusalemHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	jerusalem := location.GetJerusalem()

	loc, err := time.LoadLocation(jerusalem.Timezone)
	if err != nil {
		http.Error(w, "No se pudo cargar la zona horaria de Jerusalén", http.StatusInternalServerError)
		return
	}

	nowJerusalem := time.Now().In(loc)
	sunset := solar.ApproxSunsetJerusalem(nowJerusalem)
	afterSunset := biblical.IsAfterSunset(nowJerusalem, sunset)
	biblicalDate := biblical.GetBiblicalDate(nowJerusalem, sunset)

	monthStart := biblical.FindEstimatedMonthStart(nowJerusalem)
	biblicalDay := biblical.CalculateBiblicalDay(monthStart, nowJerusalem, afterSunset)
	biblicalMonth := biblical.GetCurrentBiblicalMonthAt(nowJerusalem)

	currentFeasts := feasts.GetCurrentFeasts(monthStart, nowJerusalem)
	upcomingFeasts := feasts.GetUpcomingFeasts(monthStart, nowJerusalem)
	omerInfo := feasts.CalculateOmerInfo(monthStart, nowJerusalem)

	data := models.FeastsResponse{
		CivilDate:      nowJerusalem.Format("2006-01-02"),
		BiblicalDate:   biblicalDate.Format("2006-01-02"),
		BiblicalMonth:  biblicalMonth,
		BiblicalDay:    biblicalDay,
		DayNote:        "Las fiestas y los días bíblicos comienzan al atardecer y continúan hasta el siguiente atardecer.",
		CurrentFeasts:  currentFeasts,
		UpcomingFeasts: upcomingFeasts,
		Omer:           omerInfo,
	}

	response := models.ApiResponse{
		Success: true,
		Message: "Fiestas bíblicas y conteo del Omer calculados correctamente",
		Data:    data,
	}

	json.NewEncoder(w).Encode(response)
}

func AvivStatusHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	status, err := aviv.RefreshStatus()
	if err != nil {
		http.Error(w, "No se pudo cargar el estado de Aviv", http.StatusInternalServerError)
		return
	}

	response := models.ApiResponse{
		Success: true,
		Message: "Estado de Aviv cargado correctamente",
		Data:    status,
	}

	json.NewEncoder(w).Encode(response)
}

func AvivReportsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	reports, err := aviv.LoadReports()
	if err != nil {
		http.Error(w, "No se pudieron cargar los reportes de Aviv", http.StatusInternalServerError)
		return
	}

	response := models.ApiResponse{
		Success: true,
		Message: "Reportes de Aviv cargados correctamente",
		Data:    reports,
	}

	json.NewEncoder(w).Encode(response)
}

func AvivReportHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodPost {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	var req models.AvivReportRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "JSON inválido", http.StatusBadRequest)
		return
	}

	if req.UserName == "" {
		http.Error(w, "El nombre de usuario es obligatorio", http.StatusBadRequest)
		return
	}

	report, err := aviv.AddReport(req.UserName, req.Confirmed, req.Notes, req.EvidenceURL)
	if err != nil {
		http.Error(w, "No se pudo registrar el reporte", http.StatusInternalServerError)
		return
	}

	response := models.ApiResponse{
		Success: true,
		Message: "Reporte de Aviv registrado correctamente",
		Data:    report,
	}

	json.NewEncoder(w).Encode(response)
}

func AvivConfirmHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodPost {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	var req models.AdminActionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "JSON inválido", http.StatusBadRequest)
		return
	}

	if req.AdminName == "" {
		http.Error(w, "El nombre del administrador es obligatorio", http.StatusBadRequest)
		return
	}

	status, err := aviv.ConfirmOfficially(req.AdminName)
	if err != nil {
		http.Error(w, "No se pudo confirmar Aviv oficialmente", http.StatusInternalServerError)
		return
	}

	response := models.ApiResponse{
		Success: true,
		Message: "Aviv confirmado oficialmente",
		Data:    status,
	}

	json.NewEncoder(w).Encode(response)
}

func AvivRejectHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodPost {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	var req models.AdminActionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "JSON inválido", http.StatusBadRequest)
		return
	}

	if req.AdminName == "" {
		http.Error(w, "El nombre del administrador es obligatorio", http.StatusBadRequest)
		return
	}

	status, err := aviv.RejectOfficially(req.AdminName)
	if err != nil {
		http.Error(w, "No se pudo rechazar Aviv oficialmente", http.StatusInternalServerError)
		return
	}

	response := models.ApiResponse{
		Success: true,
		Message: "Aviv rechazado oficialmente",
		Data:    status,
	}

	json.NewEncoder(w).Encode(response)
}

func AvivResetHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodPost {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	status, err := aviv.ResetOfficialStatus()
	if err != nil {
		http.Error(w, "No se pudo reiniciar el estado de Aviv", http.StatusInternalServerError)
		return
	}

	response := models.ApiResponse{
		Success: true,
		Message: "Estado de Aviv reiniciado correctamente",
		Data:    status,
	}

	json.NewEncoder(w).Encode(response)
}
