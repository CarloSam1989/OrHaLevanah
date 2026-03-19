package models

type Location struct {
	Name      string  `json:"name"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Timezone  string  `json:"timezone"`
}

type MoonInfo struct {
	AgeDays      float64 `json:"age_days"`
	Illumination float64 `json:"illumination"`
	PhaseName    string  `json:"phase_name"`
	LunarDay     int     `json:"lunar_day"`
}

type MoonTodayResponse struct {
	Date     string   `json:"date"`
	Location Location `json:"location"`
	Moon     MoonInfo `json:"moon"`
}

type ObservationalInfo struct {
	PossibleNewMonth bool   `json:"possible_new_month"`
	VisibilityNote   string `json:"visibility_note"`
}

type BiblicalTodayResponse struct {
	CivilDate     string            `json:"civil_date"`
	BiblicalDate  string            `json:"biblical_date"`
	JerusalemTime string            `json:"jerusalem_time"`
	SunsetTime    string            `json:"sunset_time"`
	AfterSunset   bool              `json:"after_sunset"`
	Location      Location          `json:"location"`
	Moon          MoonInfo          `json:"moon"`
	Observational ObservationalInfo `json:"observational"`
}

type BiblicalMonthResponse struct {
	CivilDate        string   `json:"civil_date"`
	BiblicalDate     string   `json:"biblical_date"`
	JerusalemTime    string   `json:"jerusalem_time"`
	SunsetTime       string   `json:"sunset_time"`
	AfterSunset      bool     `json:"after_sunset"`
	MonthStart       string   `json:"month_start"`
	BiblicalDay      int      `json:"biblical_day"`
	IsPossibleDayOne bool     `json:"is_possible_day_one"`
	Location         Location `json:"location"`
	Moon             MoonInfo `json:"moon"`
}

type ApiResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}
