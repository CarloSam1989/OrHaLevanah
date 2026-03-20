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
	DayNote       string            `json:"day_note"`
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
	DayNote          string   `json:"day_note"`
	MonthStart       string   `json:"month_start"`
	BiblicalDay      int      `json:"biblical_day"`
	IsPossibleDayOne bool     `json:"is_possible_day_one"`
	Location         Location `json:"location"`
	Moon             MoonInfo `json:"moon"`
}

type Feast struct {
	Name            string `json:"name"`
	BiblicalMonth   int    `json:"biblical_month"`
	BiblicalDay     int    `json:"biblical_day"`
	GregorianDate   string `json:"gregorian_date"`
	Weekday         string `json:"weekday"`
	BiblicalStartAt string `json:"biblical_start_at"`
	Description     string `json:"description"`
}

type OmerInfo struct {
	FirstShabbatInMatzot string `json:"first_shabbat_in_matzot"`
	BikkurimDate         string `json:"bikkurim_date"`
	BikkurimWeekday      string `json:"bikkurim_weekday"`
	ShavuotDate          string `json:"shavuot_date"`
	ShavuotWeekday       string `json:"shavuot_weekday"`
	OmerDayToday         int    `json:"omer_day_today"`
	IsOmerCounting       bool   `json:"is_omer_counting"`
	DayNote              string `json:"day_note"`
}

type FeastsResponse struct {
	CivilDate      string   `json:"civil_date"`
	BiblicalDate   string   `json:"biblical_date"`
	BiblicalMonth  int      `json:"biblical_month"`
	BiblicalDay    int      `json:"biblical_day"`
	DayNote        string   `json:"day_note"`
	CurrentFeasts  []Feast  `json:"current_feasts"`
	UpcomingFeasts []Feast  `json:"upcoming_feasts"`
	Omer           OmerInfo `json:"omer"`
}

type AvivReportRequest struct {
	UserName    string `json:"user_name"`
	Confirmed   bool   `json:"confirmed"`
	Notes       string `json:"notes"`
	EvidenceURL string `json:"evidence_url"`
}

type AdminActionRequest struct {
	AdminName string `json:"admin_name"`
}

type ApiResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}
