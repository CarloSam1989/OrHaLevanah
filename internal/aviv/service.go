package aviv

import (
	"fmt"
	"time"
)

func GetCurrentYear() int {
	return time.Now().Year()
}

func HasSpringEquinoxPassed() bool {
	now := time.Now().UTC()
	year := now.Year()

	// Aproximación simple inicial.
	// Luego lo refinamos con cálculo/tabla astronómica real.
	equinox := time.Date(year, time.March, 20, 0, 0, 0, 0, time.UTC)

	return now.Equal(equinox) || now.After(equinox)
}

func RefreshStatus() (Status, error) {
	status, err := LoadStatus()
	if err != nil {
		return status, err
	}

	reports, err := LoadReports()
	if err != nil {
		return status, err
	}

	status.Year = GetCurrentYear()
	status.EquinoxPassed = HasSpringEquinoxPassed()
	status.ReportsCount = len(reports)
	status.ConfirmCount = 0
	status.RejectCount = 0

	for _, report := range reports {
		if report.Year != status.Year {
			continue
		}

		if report.Confirmed {
			status.ConfirmCount++
		} else {
			status.RejectCount++
		}
	}

	if status.OfficialConfirmed {
		status.Status = "confirmed"
	} else if status.OfficialRejected {
		status.Status = "rejected"
	} else if status.ReportsCount >= 3 && status.ConfirmCount > status.RejectCount {
		status.Status = "community_consensus"
	} else {
		status.Status = "pending"
	}

	err = SaveStatus(status)
	return status, err
}

func AddReport(userName string, confirmed bool, notes, evidenceURL string) (Report, error) {
	reports, err := LoadReports()
	if err != nil {
		return Report{}, err
	}

	report := Report{
		ID:          fmt.Sprintf("aviv-%d", time.Now().UnixNano()),
		Year:        GetCurrentYear(),
		UserName:    userName,
		Confirmed:   confirmed,
		Notes:       notes,
		EvidenceURL: evidenceURL,
		CreatedAt:   time.Now().Format(time.RFC3339),
	}

	reports = append(reports, report)

	if err := SaveReports(reports); err != nil {
		return Report{}, err
	}

	_, err = RefreshStatus()
	if err != nil {
		return Report{}, err
	}

	return report, nil
}

func ConfirmOfficially(adminName string) (Status, error) {
	status, err := LoadStatus()
	if err != nil {
		return status, err
	}

	status.Year = GetCurrentYear()
	status.EquinoxPassed = HasSpringEquinoxPassed()
	status.OfficialConfirmed = true
	status.OfficialRejected = false
	status.OfficialConfirmedBy = adminName
	status.OfficialDate = time.Now().Format(time.RFC3339)
	status.Status = "confirmed"

	if err := SaveStatus(status); err != nil {
		return status, err
	}

	return RefreshStatus()
}

func RejectOfficially(adminName string) (Status, error) {
	status, err := LoadStatus()
	if err != nil {
		return status, err
	}

	status.Year = GetCurrentYear()
	status.EquinoxPassed = HasSpringEquinoxPassed()
	status.OfficialConfirmed = false
	status.OfficialRejected = true
	status.OfficialConfirmedBy = adminName
	status.OfficialDate = time.Now().Format(time.RFC3339)
	status.Status = "rejected"

	if err := SaveStatus(status); err != nil {
		return status, err
	}

	return RefreshStatus()
}

func ResetOfficialStatus() (Status, error) {
	status, err := LoadStatus()
	if err != nil {
		return status, err
	}

	status.Year = GetCurrentYear()
	status.EquinoxPassed = HasSpringEquinoxPassed()
	status.OfficialConfirmed = false
	status.OfficialRejected = false
	status.OfficialConfirmedBy = ""
	status.OfficialDate = ""
	status.Status = "pending"

	if err := SaveStatus(status); err != nil {
		return status, err
	}

	return RefreshStatus()
}
