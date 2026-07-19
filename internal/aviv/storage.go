package aviv

import (
	"encoding/json"
	"os"
	"strings"
)

const (
	statusFilePath  = "data/aviv_status.json"
	reportsFilePath = "data/aviv_reports.json"
)

func statusPath() string {
	if value := strings.TrimSpace(os.Getenv("AVIV_STATUS_PATH")); value != "" {
		return value
	}
	return statusFilePath
}

func reportsPath() string {
	if value := strings.TrimSpace(os.Getenv("AVIV_REPORTS_PATH")); value != "" {
		return value
	}
	return reportsFilePath
}

func LoadStatus() (Status, error) {
	var status Status

	data, err := os.ReadFile(statusPath())
	if err != nil {
		return status, err
	}

	err = json.Unmarshal(data, &status)
	return status, err
}

func SaveStatus(status Status) error {
	data, err := json.MarshalIndent(status, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(statusPath(), data, 0644)
}

func LoadReports() ([]Report, error) {
	var reports []Report

	data, err := os.ReadFile(reportsPath())
	if err != nil {
		return reports, err
	}

	err = json.Unmarshal(data, &reports)
	return reports, err
}

func SaveReports(reports []Report) error {
	data, err := json.MarshalIndent(reports, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(reportsPath(), data, 0644)
}
