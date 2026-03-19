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

type ApiResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}
