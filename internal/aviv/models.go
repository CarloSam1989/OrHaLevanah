package aviv

type Report struct {
	ID          string `json:"id"`
	Year        int    `json:"year"`
	UserName    string `json:"user_name"`
	Confirmed   bool   `json:"confirmed"`
	Notes       string `json:"notes"`
	EvidenceURL string `json:"evidence_url,omitempty"`
	CreatedAt   string `json:"created_at"`
}

type Status struct {
	Year                int    `json:"year"`
	EquinoxPassed       bool   `json:"equinox_passed"`
	OfficialConfirmed   bool   `json:"official_confirmed"`
	OfficialRejected    bool   `json:"official_rejected"`
	OfficialConfirmedBy string `json:"official_confirmed_by,omitempty"`
	OfficialDate        string `json:"official_date,omitempty"`
	ReportsCount        int    `json:"reports_count"`
	ConfirmCount        int    `json:"confirm_count"`
	RejectCount         int    `json:"reject_count"`
	Status              string `json:"status"` // pending, community_consensus, confirmed, rejected
}
