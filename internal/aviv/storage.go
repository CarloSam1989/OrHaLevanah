package aviv

import (
	"encoding/json"
	"os"
)

const (
	statusFilePath  = "data/aviv_status.json"
	reportsFilePath = "data/aviv_reports.json"
)

func LoadStatus() (Status, error) {
	var status Status

	data, err := os.ReadFile(statusFilePath)
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

	return os.WriteFile(statusFilePath, data, 0644)
}

func LoadReports() ([]Report, error) {
	var reports []Report

	data, err := os.ReadFile(reportsFilePath)
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

	return os.WriteFile(reportsFilePath, data, 0644)
}
